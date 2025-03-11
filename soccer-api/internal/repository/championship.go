package repository

import (
	"context"

	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
)

type Championship interface {
	FindAll(ctx context.Context) ([]*entity.Championship, error)
}
