package sender

import (
	"sync"

	"github.com/pkg/errors"
)

type Service struct {
	opts Options

	stopWg sync.WaitGroup
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid configuration")
	}

	return &Service{
		opts: opts,
	}, nil
}

func (s *Service) Wait() {
	s.stopWg.Wait()
}
