package queue

import (
	"encoding/json"
	"fmt"
	"time"
	"twitch_chat_analysis/internal/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
	user      string
	password  string
	host      string
	port      int
}

func New(
	queueName string,
	user string,
	password string,
	host string,
	port int,
) *Client {
	return &Client{
		queueName: queueName,
		user:      user,
		password:  password,
		host:      host,
		port:      port,
	}
}

func (c *Client) Connect(timeoutSec int) error {

	var elapsed int

	var err error
	var conn *amqp.Connection

	for elapsed <= timeoutSec {
		conn, err = amqp.Dial(
			fmt.Sprintf(
				"amqp://%s:%s@%s:%d",
				c.user,
				c.password,
				c.host,
				c.port,
			),
		)
		if err != nil {
			fmt.Println("failed to establish rabbitmq connection, retrying")
		} else {
			break
		}
		time.Sleep(time.Second)
		timeoutSec++
	}

	if conn == nil {
		return fmt.Errorf("failed to establish rabbitmq connection before timeout: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to establish rabbitmq channel: %w", err)
	}

	c.conn = conn
	c.channel = ch

	return nil
}

func (c *Client) Disconnect() error {

	if c.channel != nil && !c.channel.IsClosed() {
		if err := c.channel.Close(); err != nil {
			return fmt.Errorf("failed to close rabbitmq channel: %w", err)
		}
	}

	if c.conn != nil && !c.conn.IsClosed() {
		if err := c.conn.Close(); err != nil {
			return fmt.Errorf("failed to close rabbitmq connection: %w", err)
		}
	}

	return nil
}

func (c *Client) EnsureQueue() error {
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

func (c *Client) PublishMessage(msg domain.Message) error {

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
