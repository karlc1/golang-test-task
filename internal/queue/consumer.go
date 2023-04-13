package queue

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	*Client
	exchange string
}

func NewConsumer(
	queueName string,
	user string,
	password string,
	host string,
	port int,
	exchange string,
) *Consumer {
	return &Consumer{
		Client: &Client{
			queueName: queueName,
			user:      user,
			password:  password,
			host:      host,
			port:      port,
		},
		exchange: exchange,
	}
}

func (c *Consumer) Ensure() error {
	_, err := c.channel.QueueDeclare(
		c.queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to ensure rabbitmq queue: %w", err)
	}

	err = c.channel.ExchangeDeclare(
		c.exchange,
		"direct",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to ensure rabbitmq exhange: %w", err)
	}

	if err = c.channel.QueueBind(
		c.queueName,
		"worker",
		c.exchange,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	return nil
}

func (c *Consumer) Consume() (<-chan amqp.Delivery, error) {
	deliveries, err := c.channel.Consume(
		c.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error while consuming: %w", err)
	}
	return deliveries, nil
}
