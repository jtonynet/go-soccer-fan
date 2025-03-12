package gormrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/database"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/model"
	"gorm.io/gorm"
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
