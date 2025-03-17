package main

import (
	"encoding/json"
	"log"

	"github.com/jtonynet/go-soccer-fan/soccer-api/config"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/dto"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/pubsub"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/service"
	"github.com/streadway/amqp"
)

func main() {
	cfg := config.LoadConfig()

	mailService := service.NewMail(cfg.MailNotification)

	pubSub, err := pubsub.NewRabbitMQ(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer pubSub.Close()

	err = pubSub.Subscribe(cfg.RabbitMQ.FanNotificationsQueue, func(d amqp.Delivery) {
		var fReq dto.FanNotification
		err := json.Unmarshal(d.Body, &fReq)
		if err != nil {
			log.Printf("error deserializing JSON: %v\n", err)

			if nackErr := d.Nack(false, true); nackErr != nil {
				log.Printf("error nacking message: %v\n", nackErr)
			}
			return
		}

		mailService.SendMail(
			fReq.FanEmail,
			fReq.Title,
			fReq.Message,
		)

		if ackErr := d.Ack(false); ackErr != nil {
			log.Printf("error acking message: %v\n", ackErr)
		}
	})
	if err != nil {
		log.Fatalf("failed to subscribe to queue: %v", err)
	}

	select {}
}
