package repository

import (
	"context"

	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
)

type Team interface {
	CreateOrUpdateInBatch(ctx context.Context, tEntities []*entity.Team) ([]*entity.Team, error)
}
