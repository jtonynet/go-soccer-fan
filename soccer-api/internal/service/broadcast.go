package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/dto"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/rabbitmq"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/repository"
)

type Broadcast struct {
	pubSub   *rabbitmq.RabbitMQ
	teamRepo repository.Team
	queue    string
}

func NewBroadcast(pubsub *rabbitmq.RabbitMQ, teamRepo repository.Team, queue string) *Broadcast {
	return &Broadcast{pubsub, teamRepo, queue}
}

func (b *Broadcast) Publish(bReq *dto.BroadcastSendRequest) (*dto.BroadcastResponse, error) {
	_, err := b.teamRepo.FindFansByTeamName(context.Background(), bReq.Team)
	if err != nil {
		return nil, fmt.Errorf("team not found: %s", bReq.Team)
	}

	bReqStr, err := json.Marshal(bReq)
	if err != nil {
		log.Printf("error converting struct to string: %v\n", err)
		return nil, err
	}

	b.pubSub.Publish(b.queue, string(bReqStr))

	return nil, nil
}
