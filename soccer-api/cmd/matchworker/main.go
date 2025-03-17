package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jtonynet/go-soccer-fan/soccer-api/config"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/database"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/dto"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/rabbitmq"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/repository/gormrepo"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/service"
	"github.com/streadway/amqp"
)

func main() {
	cfg := config.LoadConfig()

	gormConn := database.NewGormCom(cfg.Database)
	teamRepo := gormrepo.NewTeam(gormConn)
	mailService := service.NewMail(cfg.MailNotification)

	rabbitMQ, err := rabbitmq.NewRabbitMQ(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	err = rabbitMQ.Subscribe(cfg.RabbitMQ.MatchNoticationsQueue, func(d amqp.Delivery) {
		var bReq dto.BroadcastSendRequest
		err := json.Unmarshal(d.Body, &bReq)
		if err != nil {
			log.Printf("error deserializing JSON: %v\n", err)

			if nackErr := d.Nack(false, true); nackErr != nil {
				log.Printf("error nacking message: %v\n", nackErr)
			}
			return
		}

		fanEntities, _ := teamRepo.FindFansByTeamName(context.Background(), bReq.Team)
		for _, fan := range fanEntities {
			mailService.SendMail(
				fan.Email,
				fmt.Sprintf("%s da partida do: %s", bReq.Type, bReq.Team),
				bReq.Message,
			)
		}

		if ackErr := d.Ack(false); ackErr != nil {
			log.Printf("error acking message: %v\n", ackErr)
		}
	})
	if err != nil {
		log.Fatalf("failed to subscribe to queue: %v", err)
	}

	select {}
}
