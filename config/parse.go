package config

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// GetConfig Получение конфигурации сервиса
func GetConfig(filename string) (*Config, error) {
	var cfg Config

	if _, err := toml.DecodeFile(filename, &cfg); err != nil {
		return nil, errors.Wrap(err, "cannot parse config")
	}

	return &cfg, nil
}
