package gormrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/database"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Competition struct {
	db *gorm.DB
}

func NewCompetition(gConn *database.GormConn) *Competition {
	return &Competition{
		db: gConn.GetDB(),
	}
}

func (c *Competition) FindAll(ctx context.Context) ([]*entity.Competition, error) {
	var cModelList []*model.Competition

	if err := c.db.WithContext(ctx).Find(&cModelList).Error; err != nil {
		return nil, err
	}

	var entityList []*entity.Competition
	for _, cModel := range cModelList {
		entityList = append(entityList, &entity.Competition{
			UID:    cModel.UID,
			Name:   cModel.Name,
			Season: cModel.Season,
		})
	}

	return entityList, nil
}

func (c *Competition) FindMatchsByCompetitionUID(ctx context.Context, uid uuid.UUID) ([]*entity.Match, error) {
	var competition model.Competition
	if err := c.db.WithContext(ctx).Where("uid = ?", uid).First(&competition).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}

		return []*entity.Match{}, nil
	}

	var mModelList []*model.Match
	if err := c.db.WithContext(ctx).
		Where("competition_id = ?", competition.ID).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Find(&mModelList).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}

		return []*entity.Match{}, nil
	}

	var entityList []*entity.Match
	for _, mModel := range mModelList {
		entityList = append(entityList, &entity.Match{
			UID:           mModel.UID,
			Round:         mModel.Round,
			CompetitionID: mModel.CompetitionID,
			HomeTeam:      &entity.Team{UID: mModel.HomeTeam.UID, Name: mModel.HomeTeam.Name},
			AwayTeam:      &entity.Team{UID: mModel.AwayTeam.UID, Name: mModel.AwayTeam.Name},
			HomeTeamScore: mModel.HomeTeamScore,
			AwayTeamScore: mModel.AwayTeamScore,
		})
	}

	return entityList, nil
}

func (c *Competition) CreateOrUpdateInBatch(ctx context.Context, cEntities []*entity.Competition) ([]*entity.Competition, error) {
	if len(cEntities) == 0 {
		return nil, nil
	}

	var cModels []model.Competition
	for _, ce := range cEntities {
		cModels = append(cModels, model.Competition{
			UID:        ce.UID,
			ExternalId: ce.ExternalId,
			Name:       ce.Name,
			Season:     ce.Season,
		})
	}

	err := c.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "external_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "season"}),
	}).Create(&cModels).Error
	if err != nil {
		return nil, err
	}

	var result []*entity.Competition
	for _, cModel := range cModels {
		result = append(result, &entity.Competition{
			ID:         cModel.ID,
			UID:        cModel.UID,
			ExternalId: cModel.ExternalId,
			Name:       cModel.Name,
			Season:     cModel.Season,
		})
	}

	return result, nil
}
