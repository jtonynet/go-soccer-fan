package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
)

type User interface {
	FindByUID(ctx context.Context, uid uuid.UUID) (*entity.User, error)
	FindByUserName(ctx context.Context, userName string) (*entity.User, error)
	Login(ctx context.Context, userName, password string) (string, error)
	Create(ctx context.Context, uEntity *entity.User) (*entity.User, error)
}
