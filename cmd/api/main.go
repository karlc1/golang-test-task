package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"twitch_chat_analysis/internal/domain"
	"twitch_chat_analysis/internal/environment"
	"twitch_chat_analysis/internal/queue"

	"github.com/gin-gonic/gin"
)

func main() {

	env := environment.MustGetApiEnv()

	queueClient := queue.NewProducer(
		env.RMQqueueName,
		env.RMQuser,
		env.RMQpassword,
		env.RMQhost,
		env.RMQport,
	)

	if err := queueClient.Connect(60); err != nil {
		log.Fatalf("failed to connect to queue client: %v", err)
	}

	if err := queueClient.EnsureQueue(); err != nil {
		log.Fatalf("failed to ensure rabbitmq queue: %v", err)
	}

	r := gin.Default()

	r.POST("/message", func(c *gin.Context) {

		input := domain.Message{}
		err := c.BindJSON(&input)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		fmt.Printf("received message:: %#+v", input)

		if err := queueClient.PublishMessage(input); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(200)
	})

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if err := queueClient.Disconnect(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	r.Run(fmt.Sprintf(":%d", env.ApiPort))
}
