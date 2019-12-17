package kafka

import (
	"github.com/pkg/errors"
)

type Consumer struct {
	opts Options
}

func NewConsumer(opts Options) (*Consumer, error) {
	if err := opts.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid configuration")
	}

	return &Consumer{
		opts: opts,
	}, nil
}
