package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"twitch_chat_analysis/internal/domain"
	"twitch_chat_analysis/internal/environment"
	"twitch_chat_analysis/internal/queue"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
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

		c.JSON(200, "worked3")
	})

	r.Run(fmt.Sprintf(":%d", env.ApiPort))
}

// Here we set the way error messages are displayed in the terminal.
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func rmqTest2() {
	// Here we connect to RabbitMQ or send a message if there are any errors connecting.
	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	// conn, err := amqp.Dial("amqp://user:password@localhost:7001/")
	conn, err := amqp.Dial("amqp://user:password@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// We create a Queue to send the message to.
	q, err := ch.QueueDeclare(
		"golang-queue", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// We set the payload for the message.
	body := "Golang is awesome - Keep Moving Forward!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	// If there is an error publishing the message, a log will be displayed in the terminal.
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Congrats, sending message: %s", body)
}

func rmqTest() {
	conn, err := amqp.Dial("amqp://user:password@localhost:7000/")
	if err != nil {
		panic("dial: " + err.Error())
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
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
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		panic(err)
	}

	log.Printf(" [x] Sent %s\n", body)
}
