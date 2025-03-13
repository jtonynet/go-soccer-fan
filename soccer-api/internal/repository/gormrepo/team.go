package gormrepo

import (
	"context"

	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/database"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Team struct {
	db *gorm.DB
}

func NewTeam(gConn *database.GormConn) *Team {
	return &Team{
		db: gConn.GetDB(),
	}
}

func (t *Team) CreateOrUpdateInBatch(ctx context.Context, aEntities []*entity.Team) error {
	if len(aEntities) == 0 {
		return nil
	}

	var aModels []model.Team
	for _, e := range aEntities {
		aModels = append(aModels, model.Team{
			UID:        e.UID,
			ExternalId: e.ExternalId,
			Name:       e.Name,
			FullName:   e.FullName,
		})
	}

	err := t.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "external_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "full_name"}),
	}).Create(&aModels).Error

	return err
}
