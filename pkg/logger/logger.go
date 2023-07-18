package logger

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type rabbitMQConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Exchange string
}

type RabbitMQLogger struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	config     rabbitMQConfig
}

func buildEnvRabbitMQConfig() rabbitMQConfig {
	return rabbitMQConfig{
		Host:     "rabbitmq",
		Port:     "5672",
		Username: "guest",
		Password: "guest",
		Exchange: "logs",
	}
}

func dialWithRabbitMQ(cfg *rabbitMQConfig) *amqp.Connection {
	rabbitURL := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	connection, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatal(err, "unable connect to RabbitMQ")
		return nil
	}

	return connection
}

func configureRabbitMQChannel(cfg *rabbitMQConfig, connection *amqp.Connection) *amqp.Channel {
	channel, err := connection.Channel()
	if err != nil {
		log.Fatal(err, "unable to configure RabbitMQ channel")
		return nil
	}

	if err := channel.ExchangeDeclare(
		cfg.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Fatal(err, "unable to declare exchange")
	}
	return channel
}

func NewRabbitMQLogger() *RabbitMQLogger {

	cfg := buildEnvRabbitMQConfig()
	connection := dialWithRabbitMQ(&cfg)
	channel := configureRabbitMQChannel(&cfg, connection)

	logger := &RabbitMQLogger{
		config:     cfg,
		channel:    channel,
		connection: connection,
	}

	return logger
}

func (l *RabbitMQLogger) publishLog(key string, message string) {
	err := l.channel.PublishWithContext(context.Background(),
		l.config.Exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Fatal(err.Error(), "failed to publish a message")
	}
}

func (l *RabbitMQLogger) Error(message string) {
	log.Println(message)
	l.publishLog("error", message)
}

func (l *RabbitMQLogger) Info(message string) {
	log.Println(message)
	l.publishLog("info", message)
}
