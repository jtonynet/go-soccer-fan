package repository

import (
	"context"

	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
)

type Fan interface {
	Create(ctx context.Context, fEntity *entity.Fan) (*entity.Fan, error)
}
