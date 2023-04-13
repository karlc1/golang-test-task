package queue

import (
	"encoding/json"
	"fmt"
	"twitch_chat_analysis/internal/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	*Client
}

func NewProducer(
	queueName string,
	user string,
	password string,
	host string,
	port int,
) *Producer {
	return &Producer{
		Client: &Client{
			queueName: queueName,
			user:      user,
			password:  password,
			host:      host,
			port:      port,
		},
	}
}

func (c *Producer) EnsureQueue() error {
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
	return nil
}

func (c *Producer) PublishMessage(msg domain.Message) error {

	b, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message for publish: %w", err)
	}

	err = c.channel.Publish(
		"",
		c.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		})

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}
