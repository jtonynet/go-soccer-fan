package rabbitmq

import (
	"fmt"
	"log"
	"os"

	"github.com/jtonynet/go-soccer-fan/soccer-api/config"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ(_ *config.RabbitMQ) (*RabbitMQ, error) {
	user := os.Getenv("RABBITMQ_USER")
	password := os.Getenv("RABBITMQ_PASS")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")

	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)

	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	rabbitMQ := &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}

	return rabbitMQ, nil
}

func (r *RabbitMQ) Publish(queueName string, body string) error {
	q, err := r.Channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	err = r.Channel.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	log.Printf(" [x] Sent %s", body)
	return nil
}

func (r *RabbitMQ) Close() {
	r.Channel.Close()
	r.Conn.Close()
}
