package gormrepo

import (
	"context"
	"fmt"

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
		return nil, nil
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

func (t *Team) FindByTeamName(ctx context.Context, tName string) (*entity.Team, error) {
	var tModel *model.Team
	if err := t.db.WithContext(ctx).Where("name = ?", tName).First(&tModel).Error; err != nil {
		return nil, fmt.Errorf("team not found: %s", tName)
	}

	return &entity.Team{
		ID:         tModel.ID,
		UID:        tModel.UID,
		ExternalID: tModel.ExternalId,
		Name:       tModel.Name,
		FullName:   tModel.FullName,
	}, nil
}

func (t *Team) FindFansByTeamName(ctx context.Context, tName string) ([]*entity.Fan, error) {
	var tModel *model.Team
	if err := t.db.WithContext(ctx).Where("name = ?", tName).First(&tModel).Error; err != nil {
		return nil, fmt.Errorf("team not found: %s", tName)
	}

	var fModels []*model.Fan
	if err := t.db.WithContext(ctx).Where("team_id = ?", tModel.ID).Find(&fModels).Error; err != nil {
		return nil, fmt.Errorf("failed to find fans: %w", err)
	}

	var fEntities []*entity.Fan
	for _, fModel := range fModels {
		fEntities = append(fEntities, &entity.Fan{
			UID:   fModel.UID,
			Name:  fModel.Name,
			Email: fModel.Email,
			Team: &entity.Team{
				ID:         tModel.ID,
				UID:        tModel.UID,
				ExternalID: tModel.ExternalId,
				Name:       tModel.Name,
				FullName:   tModel.FullName,
			},
		})
	}

	return fEntities, nil
}
