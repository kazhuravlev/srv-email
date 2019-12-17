package sender

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

func _smtpOptValidator(o *Options) error {
	if utils.IsNil(o.smtp) {
		return errors.Wrap(ErrInvalidOption, "Smtp must be set (type Smtp)")
	}
	return nil
}

func _consumerOptValidator(o *Options) error {
	if utils.IsNil(o.consumer) {
		return errors.Wrap(ErrInvalidOption, "Consumer must be set (type Consumer)")
	}
	return nil
}

func _topicOptValidator(o *Options) error {
	if utils.IsNil(o.topic) {
		return errors.Wrap(ErrInvalidOption, "Topic must be set (type string)")
	}
	return nil
}

func _unmarshalerOptValidator(o *Options) error {
	if utils.IsNil(o.unmarshaler) {
		return errors.Wrap(ErrInvalidOption, "Unmarshaler must be set (type MessageUnmarshaler)")
	}
	return nil
}

func _maxWorkersOptValidator(o *Options) error {
	if utils.IsNil(o.maxWorkers) {
		return errors.Wrap(ErrInvalidOption, "MaxWorkers must be set (type int)")
	}
	return nil
}

func NewOptions(
	logger Logger,
	smtp Smtp,
	consumer Consumer,
	topic string,
	unmarshaler MessageUnmarshaler,
	maxWorkers int,

	options ...optMeta,
) Options {
	o := Options{}
	o.logger = logger
	o.smtp = smtp
	o.consumer = consumer
	o.topic = topic
	o.unmarshaler = unmarshaler
	o.maxWorkers = maxWorkers

	for i := range options {
		options[i].setter(&o)
	}

	return o
}

func (o *Options) Validate() error {
	if err := _loggerOptValidator(o); err != nil {
		return errors.Wrap(err, "invalid value for option WithLogger")
	}
	if err := _smtpOptValidator(o); err != nil {
		return errors.Wrap(err, "invalid value for option WithSmtp")
	}
	if err := _consumerOptValidator(o); err != nil {
		return errors.Wrap(err, "invalid value for option WithConsumer")
	}
	if err := _topicOptValidator(o); err != nil {
		return errors.Wrap(err, "invalid value for option WithTopic")
	}
	if err := _unmarshalerOptValidator(o); err != nil {
		return errors.Wrap(err, "invalid value for option WithUnmarshaler")
	}
	if err := _maxWorkersOptValidator(o); err != nil {
		return errors.Wrap(err, "invalid value for option WithMaxWorkers")
	}

	return nil
}
