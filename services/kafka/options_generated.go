package kafka

import (
	"github.com/kazhuravlev/options-gen/generator/utils"
	"github.com/pkg/errors"
)

type optMeta struct {
	setter    func(o *Options)
	validator func(o *Options) error
}

func _loggerOptValidator(o *Options) error {
	if utils.IsNil(o.logger) {
		return errors.Wrap(ErrInvalidOption, "Logger must be set (type Logger)")
	}
	return nil
}

func _brokersOptValidator(o *Options) error {
	if utils.IsNil(o.brokers) {
		return errors.Wrap(ErrInvalidOption, "Brokers must be set (type []string)")
	}
	return nil
}

func _clientIDOptValidator(o *Options) error {
	if utils.IsNil(o.clientID) {
		return errors.Wrap(ErrInvalidOption, "ClientID must be set (type string)")
	}
	return nil
}

func NewOptions(
	logger Logger,
	brokers []string,
	clientID string,

	options ...optMeta,
) Options {
	o := Options{}
	o.logger = logger
	o.brokers = brokers
	o.clientID = clientID

	for i := range options {
		options[i].setter(&o)
	}

	return o
}

func (o *Options) Validate() error {
	if err := _loggerOptValidator(o); err != nil {
		return errors.Wrap(err, "invalid value for option WithLogger")
	}
	if err := _brokersOptValidator(o); err != nil {
		return errors.Wrap(err, "invalid value for option WithBrokers")
	}
	if err := _clientIDOptValidator(o); err != nil {
		return errors.Wrap(err, "invalid value for option WithClientID")
	}

	return nil
}
