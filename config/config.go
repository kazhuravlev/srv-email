package config

import "strings"

type SenderConfig struct {
	MaxWorkers  int    `toml:"max_workers"`
	Unmarshaler string `toml:"unmarshaler"`
}

type SmtpConfig struct {
	ServerAddr string `toml:"server_addr"`
	Username   string `toml:"username"`
	Password   string `toml:"password"`
}

type KafkaConfig struct {
	Brokers  string `toml:"brokers"`
	ClientID string `toml:"client_id"`
	Topic    string `toml:"topic"`
}

func (c *KafkaConfig) GetBrokers() []string {
	return strings.Split(c.Brokers, ",")
}

type Config struct {
	Sender SenderConfig `toml:"sender"`
	Smtp   SmtpConfig   `toml:"smtp"`
	Kafka  KafkaConfig  `toml:"kafka"`
}
