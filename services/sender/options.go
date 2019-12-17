package sender

import (
	"github.com/pkg/errors"
)

var ErrInvalidOption = errors.New("invalid option")

//go:generate options-gen -filename=$GOFILE -out-filename=options_generated.go -pkg=sender -from-struct=Options
type Options struct {
	logger      Logger             `option:"required,not-empty"`
	smtp        Smtp               `option:"required,not-empty"`
	consumer    Consumer           `option:"required,not-empty"`
	topic       string             `option:"required,not-empty"`
	unmarshaler MessageUnmarshaler `option:"required,not-empty"`
	maxWorkers  int                `option:"required,not-empty"`
}
