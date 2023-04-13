package cache

import (
	"context"
	"fmt"
	"twitch_chat_analysis/internal/domain"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	redis *redis.Client
}

func New(host string, port int) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
	})

	return &Client{
		redis: rdb,
	}
}

func (c *Client) Put(ctx context.Context, msg domain.Message) error {
	return c.redis.RPush(ctx, msgToKey(msg), msg.Message).Err()
}

func msgToKey(msg domain.Message) string {
	return fmt.Sprintf("%s-%s", msg.Sender, msg.Receiver)
}
