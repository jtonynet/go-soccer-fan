package gormrepo

import (
	"context"

	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/database"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/model"
	"gorm.io/gorm"
)

type Fan struct {
	db *gorm.DB
}

func NewFan(gConn *database.GormConn) *Fan {
	return &Fan{
		db: gConn.GetDB(),
	}
}

func (f *Fan) Create(ctx context.Context, fEntity *entity.Fan) (*entity.Fan, error) {

	var team model.Team
	if err := f.db.WithContext(ctx).Where("name = ?", fEntity.Team.Name).First(&team).Error; err != nil {
		return nil, err
	}

	fModel := &model.Fan{
		UID:    fEntity.UID,
		Name:   fEntity.Name,
		Email:  fEntity.Email,
		TeamID: int(team.ID),
	}

	if err := f.db.WithContext(ctx).Create(fModel).Error; err != nil {
		return nil, err
	}

	return &entity.Fan{
		UID:   fModel.UID,
		Name:  fModel.Name,
		Email: fModel.Email,
		Team: &entity.Team{
			ID:   int(team.ID),
			UID:  team.UID,
			Name: team.Name,
		},
	}, nil
}
