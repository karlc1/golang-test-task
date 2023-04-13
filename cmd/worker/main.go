package main

import (
	"context"
	"encoding/json"
	"fmt"
	"twitch_chat_analysis/internal/cache"
	"twitch_chat_analysis/internal/domain"
	"twitch_chat_analysis/internal/environment"
	"twitch_chat_analysis/internal/queue"
)

func main() {
	env := environment.MustGetWorkerEnv()

	cacheClient := cache.New(env.RedisHost, env.RedisPort)

	queueClient := queue.NewConsumer(
		env.RMQqueueName,
		env.RMQuser,
		env.RMQpassword,
		env.RMQhost,
		env.RMQport,
		env.RMQexchange,
	)

	err := queueClient.Connect(60)
	if err != nil {
		panic(fmt.Sprintf("failed to connect within time limit: %v", err))
	}

	err = queueClient.Ensure()
	if err != nil {
		panic(fmt.Sprintf("failed to ensure queue connection: %v", err))
	}

	deliveries, err := queueClient.Consume()
	if err != nil {
		panic(fmt.Sprintf("failed to start consuming from queue: %v", err))
	}

	for d := range deliveries {

		fmt.Printf("received message: %s\n", string(d.Body))

		msg := domain.Message{}
		if err := json.Unmarshal(d.Body, &msg); err != nil {
			fmt.Printf("failed to decode malformed body: %s\n", string(d.Body))
			d.Nack(false, false)
			continue
		}

		if err := cacheClient.Put(context.Background(), msg); err != nil {
			fmt.Printf("failed to put msg on redis: %v\n", err)
			d.Nack(false, true)
			continue
		}

		fmt.Println("message handle success")
		d.Ack(false)
	}
}
