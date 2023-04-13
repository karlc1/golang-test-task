package queue

import (
	"fmt"
	"time"

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
