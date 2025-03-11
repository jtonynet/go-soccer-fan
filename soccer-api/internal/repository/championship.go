package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
)

type Championship interface {
	FindAll(ctx context.Context) ([]*entity.Championship, error)
	FindMatchsByChampionshipUID(ctx context.Context, uid uuid.UUID) ([]*entity.Match, error)
}
