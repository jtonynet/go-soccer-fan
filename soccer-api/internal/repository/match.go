package repository

import (
	"context"

	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
)

type Match interface {
	CreateOrUpdateInBatch(ctx context.Context, mEntities []*entity.Match) ([]*entity.Match, error)
}
