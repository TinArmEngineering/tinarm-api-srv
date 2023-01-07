package dbo

import (
	"context"
	"log"

	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// Dev config
	QUEUE_HOST = "localhost:5672"
)

func qHost() string {
	return getEnvironment("GOSERVER_QUEUE_HOST", QUEUE_HOST)
}

func queueConnectionString() string {
	return "amqp://guest:guest@" + qHost() + "/"
}

// Post to RabbitMQ
func Enqueue(body string) {

	log.Printf("Queueing: " + body)

	conn, err := amqp.Dial(queueConnectionString())
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rectangle_mesh_queue", // name
		true,                   // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	failOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}
