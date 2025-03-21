package repository

import (
	"context"

	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
)

type Team interface {
	FindByTeamName(ctx context.Context, tName string) (*entity.Team, error)
	FindFansByTeamName(ctx context.Context, tName string) ([]*entity.Fan, error)
	CreateOrUpdateInBatch(ctx context.Context, tEntities []*entity.Team) ([]*entity.Team, error)
}
