package gormrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/database"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/model"
	"gorm.io/gorm"
)

type Championship struct {
	db *gorm.DB
}

func NewChampionship(gConn *database.GormConn) *Championship {
	return &Championship{
		db: gConn.GetDB(),
	}
}

func (c *Championship) FindAll(ctx context.Context) ([]*entity.Championship, error) {
	var cModelList []*model.Championship

	if err := c.db.WithContext(ctx).Find(&cModelList).Error; err != nil {
		return nil, err
	}

	var entityList []*entity.Championship
	for _, cModel := range cModelList {
		entityList = append(entityList, &entity.Championship{
			UID:    cModel.UID,
			Name:   cModel.Name,
			Season: cModel.Season,
		})
	}

	return entityList, nil
}

func (c *Championship) FindMatchsByChampionshipUID(ctx context.Context, uid uuid.UUID) ([]*entity.Match, error) {
	var championship model.Championship
	if err := c.db.WithContext(ctx).Where("uid = ?", uid).First(&championship).Error; err != nil {
		return nil, err
	}

	var mModelList []*model.Match
	if err := c.db.WithContext(ctx).
		Where("championship_id = ?", championship.ID).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Find(&mModelList).Error; err != nil {
		return nil, err
	}

	var entityList []*entity.Match
	for _, mModel := range mModelList {
		entityList = append(entityList, &entity.Match{
			UID:            mModel.UID,
			Round:          mModel.Round,
			ChampionshipID: mModel.ChampionshipID,
			HomeTeam:       &entity.Team{UID: mModel.HomeTeam.UID, Name: mModel.HomeTeam.Name},
			AwayTeam:       &entity.Team{UID: mModel.AwayTeam.UID, Name: mModel.AwayTeam.Name},
			HomeTeamScore:  mModel.HomeTeamScore,
			AwayTeamScore:  mModel.AwayTeamScore,
		})
	}

	return entityList, nil
}
