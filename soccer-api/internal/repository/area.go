package repository

import (
	"context"

	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
)

type Area interface {
	CreateOrUpdateInBatch(ctx context.Context, aEntities []*entity.Area) error
}
