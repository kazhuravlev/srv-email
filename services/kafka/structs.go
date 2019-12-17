package kafka

import (
	"context"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/segmentio/kafka-go"
)

const (
	headerKeyNumberOfAttempts = "x-number-of-attempts"
)

var (
	_ Message = new(message)
)

type Message interface {
	Ack(ctx context.Context) error
	Delay(ctx context.Context, retryCount int) error
	GetBody() []byte
}

type message struct {
	srcMessage kafka.Message
	reader     *kafka.Reader
	writer     *kafka.Writer

	Body []byte
}

func (m *message) GetBody() []byte {
	return m.Body
}

func (m *message) Ack(ctx context.Context) error {
	return m.reader.CommitMessages(ctx, m.srcMessage)
}

func (m *message) GetNumberOfAttempts() int {
	for i := range m.srcMessage.Headers {
		hdr := &m.srcMessage.Headers[i]
		if hdr.Key != headerKeyNumberOfAttempts {
			continue
		}

		attempts, err := strconv.Atoi(string(hdr.Value))
		if err == nil {
			return attempts
		}
		break
	}

	return 1
}

func (m *message) CopyKafkaMessage() kafka.Message {
	return kafka.Message{
		Topic:     m.srcMessage.Topic,
		Partition: 0,
		Offset:    0,
		Key:       nil,
		Value:     m.srcMessage.Value,
		Headers:   m.srcMessage.Headers,
		Time:      time.Time{},
	}
}

func (m *message) Delay(ctx context.Context, retryCount int) error {
	attempts := m.GetNumberOfAttempts()
	if attempts >= retryCount {
		return nil
	}

	msg := m.CopyKafkaMessage()
	setNumberOfAttempts(&msg, attempts+1)

	if err := m.Ack(ctx); err != nil {
		return errors.Wrap(err, "cannot ack message")
	}

	return m.writer.WriteMessages(ctx, msg)
}

func setNumberOfAttempts(msg *kafka.Message, attempts int) {
	val := []byte(strconv.Itoa(attempts))

	for i := range msg.Headers {
		hdr := &msg.Headers[i]

		if hdr.Key == headerKeyNumberOfAttempts {
			msg.Headers[i].Value = val
			return
		}
	}

	msg.Headers = append(msg.Headers, kafka.Header{
		Key:   headerKeyNumberOfAttempts,
		Value: val,
	})
}
