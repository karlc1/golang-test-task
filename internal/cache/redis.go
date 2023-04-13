package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"twitch_chat_analysis/internal/domain"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	redis *redis.Client
}

func New(host string, port int) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", host, port),
	})

	return &Client{
		redis: rdb,
	}
}

func (c *Client) Put(ctx context.Context, msg domain.Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.redis.RPush(ctx, msgToKey(msg), string(b)).Err()
}

func (c *Client) List(ctx context.Context, sender, receiver string) ([]domain.Message, error) {
	res, err := c.redis.LRange(ctx, key(sender, receiver), 0, -1).Result()
	if err != nil {
		return nil, err
	}

	ret := make([]domain.Message, len(res))
	for _, s := range res {
		m := domain.Message{}
		if err := json.Unmarshal([]byte(s), &m); err != nil {
			return nil, err
		}
		ret = append(ret, m)
	}

	return ret, nil
}

func msgToKey(msg domain.Message) string {
	return key(msg.Sender, msg.Receiver)
}

func key(sender, receiver string) string {
	return fmt.Sprintf("%s-%s", sender, receiver)
}
