package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
)

type Competition interface {
	FindAll(ctx context.Context) ([]*entity.Competition, error)
	FindMatchsByCompetitionUID(ctx context.Context, uid uuid.UUID) ([]*entity.Match, error)
}
