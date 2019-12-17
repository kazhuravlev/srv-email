package kafka

import (
	"github.com/pkg/errors"
)

var ErrInvalidOption = errors.New("invalid option")

//go:generate options-gen -filename=$GOFILE -out-filename=options_generated.go -pkg=kafka -from-struct=Options
type Options struct {
	logger   Logger   `option:"required,not-empty"`
	brokers  []string `option:"required,not-empty"`
	clientID string   `option:"required,not-empty"`
}
