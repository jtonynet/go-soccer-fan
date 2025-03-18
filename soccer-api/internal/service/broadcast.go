package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/dto"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/pubsub"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/repository"
)

type Broadcast struct {
	pubSub   *pubsub.RabbitMQ
	teamRepo repository.Team
	queue    string
}

func NewBroadcast(pubsub *pubsub.RabbitMQ, teamRepo repository.Team, queue string) *Broadcast {
	return &Broadcast{pubsub, teamRepo, queue}
}

func (b *Broadcast) Publish(bReq *dto.BroadcastSendRequest) (*dto.BroadcastResponse, error) {
	_, err := b.teamRepo.FindFansByTeamName(context.Background(), bReq.TeamName)
	if err != nil {
		return nil, errors.New("not found")
	}

	bReqStr, err := json.Marshal(bReq)
	if err != nil {
		log.Printf("error converting struct to string: %v\n", err)
		return nil, err
	}

	b.pubSub.Publish(b.queue, string(bReqStr))

	return &dto.BroadcastResponse{
		Message: "Notificação enviada",
	}, nil
}
