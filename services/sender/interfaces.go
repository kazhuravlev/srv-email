package sender

import (
	"context"

	"github.com/kazhuravlev/srv-email/services/smtp"

	"go.uber.org/zap"

	"github.com/kazhuravlev/srv-email/services/kafka"
)

type Consumer interface {
	Subscribe(ctx context.Context, topic string) (<-chan kafka.Message, error)
}

type Smtp interface {
	SendEmail(ctx context.Context, email smtp.Email) error
}

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

type MessageUnmarshaler interface {
	Unmarshal([]byte, interface{}) error
}
