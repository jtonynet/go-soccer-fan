package rabbitmq

import (
	"fmt"
	"log"

	"github.com/jtonynet/go-soccer-fan/soccer-api/config"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ(cfg *config.RabbitMQ) (*RabbitMQ, error) {
	user := cfg.User
	password := cfg.Password
	host := cfg.Host
	port := cfg.Port

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

	return nil
}

func (r *RabbitMQ) Subscribe(queueName string, handler func(amqp.Delivery)) error {
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

	msgs, err := r.Channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	go func() {
		for d := range msgs {
			handler(d)
		}
	}()

	log.Printf(" [*] Waiting for messages in %s. To exit press CTRL+C", queueName)
	select {}
}

func (r *RabbitMQ) Close() {
	r.Channel.Close()
	r.Conn.Close()
}
