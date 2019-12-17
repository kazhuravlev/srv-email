package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kazhuravlev/srv-email/config"

	"github.com/kazhuravlev/srv-email/services/sender"

	"go.uber.org/zap"

	"github.com/kazhuravlev/srv-email/services/smtp"

	"github.com/kazhuravlev/srv-email/services/kafka"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

		<-signalCh
		signal.Stop(signalCh)
		cancel()
	}()

	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("cannot init logger")
		os.Exit(1)
	}

	cfg, err := config.GetConfig("config.toml")
	if err != nil {
		logger.Fatal("cannot parse config", zap.Error(err))
	}

	consumer, err := kafka.NewConsumer(kafka.NewOptions(
		logger,
		cfg.Kafka.GetBrokers(),
		cfg.Kafka.ClientID,
	))
	if err != nil {
		logger.Fatal("cannot create new consumer", zap.Error(err))
	}

	smtpClient, err := smtp.New(cfg.Smtp.ServerAddr, cfg.Smtp.Username, cfg.Smtp.Password)
	if err != nil {
		logger.Fatal("cannot create new smtp client", zap.Error(err))
	}

	var unmarshaler sender.MessageUnmarshaler
	switch cfg.Sender.Unmarshaler {
	case "json":
		unmarshaler = sender.JSONUnmarshaler{}
	case "proto":
		unmarshaler = sender.ProtoUnmarshaler{}
	default:
		logger.Fatal("unknown unmarshaler: " + cfg.Sender.Unmarshaler)
	}

	emailSender, err := sender.New(sender.NewOptions(
		logger,
		smtpClient,
		consumer,
		cfg.Kafka.Topic,
		unmarshaler,
		cfg.Sender.MaxWorkers,
	))
	if err != nil {
		logger.Fatal("cannot create new sender", zap.Error(err))
	}

	if err := emailSender.Run(ctx); err != nil {
		logger.Fatal("cannot run sender", zap.Error(err))
	}

	emailSender.Wait()
}
