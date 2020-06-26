package email

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/streadway/amqp"
)

var (
	RabbitMQChan *amqp.Channel
	RabbitMQQueue amqp.Queue
)

type RabbitMqEnv struct {
	RabbitMqHost string `env:"RABBITMQ_HOST"`
	RabbitMqPort string `env:"RABBITMQ_PORT"`
	RabbitMqUser string `env:"RABBITMQ_DEFAULT_USER"`
	RabbitMqPass string `env:"RABBITMQ_DEFAULT_PASS"`
}

func ConnectToRabbitMQ() error {
	cfg := RabbitMqEnv{}
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("failed to parse env: %v", err)
	}

	fmt.Printf("%+v", cfg)

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMqPass,
		cfg.RabbitMqUser,
		cfg.RabbitMqHost,
		cfg.RabbitMqPort,
	))
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}

	RabbitMQChan = ch
	RabbitMQQueue = q

	return nil
}