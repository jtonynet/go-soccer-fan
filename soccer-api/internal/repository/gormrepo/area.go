package gormrepo

import (
	"context"

	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/database"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Area struct {
	db *gorm.DB
}

func NewArea(gConn *database.GormConn) *Area {
	return &Area{
		db: gConn.GetDB(),
	}
}

func (a *Area) CreateOrUpdateInBatch(ctx context.Context, aEntities []*entity.Area) error {
	if len(aEntities) == 0 {
		return nil
	}

	var aModels []model.Area
	for _, e := range aEntities {
		aModels = append(aModels, model.Area{
			ExternalId:  e.ExternalId,
			Name:        e.Name,
			CountryCode: e.CountryCode,
		})
	}

	err := a.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "external_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "country_code"}),
	}).Create(&aModels).Error

	return err
}
