package sender

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/kazhuravlev/srv-email/contracts"
)

const retryCount = 3

func (s *Service) Run(ctx context.Context) error {
	messagesCh, err := s.opts.consumer.Subscribe(ctx, s.opts.topic)
	if err != nil {
		return errors.Wrap(err, "cannot subscribe on topic")
	}

	tasksCh := make(chan taskContainer, s.opts.maxWorkers)
	s.stopWg.Add(s.opts.maxWorkers)
	for i := 0; i < s.opts.maxWorkers; i++ {
		go func() {
			defer s.stopWg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case task, ok := <-tasksCh:
					if !ok {
						return
					}

					if err := s.handleTask(ctx, task); err != nil {
						s.opts.logger.Error("cannot handle task", zap.Error(err))
						return
					}
				}
			}
		}()
	}

	s.stopWg.Add(1)
	go func() {
		defer s.stopWg.Done()
		defer close(tasksCh)

		for {
			select {
			case <-ctx.Done():
				return
			case m, ok := <-messagesCh:
				if !ok {
					return
				}

				var task contracts.Msg
				if err := s.opts.unmarshaler.Unmarshal(m.GetBody(), &task); err != nil {
					s.opts.logger.Error("cannot unmarshal message. skip",
						zap.Error(err),
						zap.Any("message", m))

					if err := m.Delay(ctx, retryCount); err != nil {
						s.opts.logger.Error("cannot delay message",
							zap.Error(err),
							zap.Any("message", m))
						return
					}
					continue
				}

				select {
				case <-ctx.Done():
				case tasksCh <- taskContainer{
					Message:      task,
					KafkaMessage: m,
				}:
				}
			}
		}
	}()

	return nil
}

func (s *Service) handleTask(ctx context.Context, task taskContainer) error {
	email := adaptTaskToEmail(task.Message)
	s.opts.logger.Debug("try to send new email...",
		zap.String("subject", email.Subject))
	if err := s.opts.smtp.SendEmail(ctx, email); err != nil {
		s.opts.logger.Info("cannot send email. try to delay task",
			zap.Error(err),
			zap.Any("email", email))

		if errDelay := task.KafkaMessage.Delay(ctx, retryCount); errDelay != nil {
			return errors.Wrap(errDelay, "cannot delay task")
		}

		return nil
	}

	s.opts.logger.Debug("try to ack task...", zap.String("subject", email.Subject))
	if err := task.KafkaMessage.Ack(ctx); err != nil {
		return errors.Wrap(err, "cannot ack message")
	}

	return nil
}
