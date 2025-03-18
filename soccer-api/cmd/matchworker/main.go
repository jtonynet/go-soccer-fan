package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jtonynet/go-soccer-fan/soccer-api/config"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/database"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/dto"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/pubsub"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/repository/gormrepo"
	"github.com/streadway/amqp"
)

func main() {
	/*
		TODO:
			Criar Handler e Service especificos para esse worker
	*/

	cfg := config.LoadConfig()

	gormConn, err := database.NewGormCom(cfg.Database)
	if err != nil {
		log.Fatalf("can't connect to database: %v", err)
	}

	teamRepo := gormrepo.NewTeam(gormConn)

	pubSub, err := pubsub.NewRabbitMQ(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer pubSub.Close()

	err = pubSub.Subscribe(cfg.RabbitMQ.MatchNoticationsQueue, func(d amqp.Delivery) {
		var bReq dto.BroadcastSendRequest
		err := json.Unmarshal(d.Body, &bReq)
		if err != nil {
			log.Printf("error deserializing JSON: %v\n", err)

			if nackErr := d.Nack(false, true); nackErr != nil {
				log.Printf("error nacking message: %v\n", nackErr)
			}
			return
		}

		score := "-"
		if bReq.Score != "" {
			score = bReq.Score
		}

		fanEntities, _ := teamRepo.FindFansByTeamName(context.Background(), bReq.TeamName)
		for _, fan := range fanEntities {
			fanNotification := &dto.FanNotification{
				FanUID:   fan.UID,
				FanEmail: fan.Email,
				Title:    fmt.Sprintf("%s da partida do: %s", bReq.Type, bReq.TeamName),
				Team:     bReq.TeamName,
				Score:    score,
				Message:  bReq.Message,
			}

			fNotify, err := json.Marshal(fanNotification)
			if err != nil {
				log.Printf("error converting struct to string: %v\n", err)
			}

			err = pubSub.Publish(cfg.RabbitMQ.FanNotificationsQueue, string(fNotify))
			if err != nil {
				log.Printf("error publishing fan notification: %v\n", err)
			}
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
