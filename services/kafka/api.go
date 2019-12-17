package kafka

import (
	"context"

	"github.com/alecthomas/units"
	"go.uber.org/zap"

	"github.com/segmentio/kafka-go"
)

func (c *Consumer) Subscribe(ctx context.Context, topic string) (<-chan Message, error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  c.opts.brokers,
		GroupID:  "srv-email-group",
		Topic:    topic,
		MinBytes: 1,
		MaxBytes: int(10 * units.MB),
	})
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:     c.opts.brokers,
		Topic:       topic,
		MaxAttempts: 3,
		BatchSize:   1,
	})

	go func() {
		<-ctx.Done()
		_ = r.Close()
	}()

	ch := make(chan Message)
	go func() {
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			m, err := r.FetchMessage(ctx)
			if err != nil {
				c.opts.logger.Error("cannot read message", zap.Error(err))
				return
			}

			ch <- &message{
				srcMessage: m,
				reader:     r,
				writer:     w,

				Body: m.Value,
			}
		}
	}()

	return ch, nil
}
