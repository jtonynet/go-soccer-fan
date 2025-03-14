package gormrepo

import (
	"context"
	"errors"

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

func (t *Team) CreateOrUpdateInBatch(ctx context.Context, aEntities []*entity.Team) ([]*entity.Team, error) {
	if len(aEntities) == 0 {
		return nil, errors.New("list of entities is empty")
	}

	var tModels []model.Team
	for _, e := range aEntities {
		tModels = append(tModels, model.Team{
			UID:        e.UID,
			ExternalId: e.ExternalID,
			Name:       e.Name,
			FullName:   e.FullName,
		})
	}

	err := t.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "external_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "full_name"}),
	}).Create(&tModels).Error
	if err != nil {
		return nil, err
	}

	var result []*entity.Team
	for _, tModel := range tModels {
		result = append(result, &entity.Team{
			ID:         tModel.ID,
			UID:        tModel.UID,
			ExternalID: tModel.ExternalId,
			Name:       tModel.Name,
			FullName:   tModel.FullName,
		})
	}

	return result, nil
}
