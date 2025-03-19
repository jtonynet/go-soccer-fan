package gormrepo

import (
	"context"

	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/database"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Match struct {
	db *gorm.DB
}

func NewMatch(gConn *database.GormConn) *Match {
	return &Match{
		db: gConn.GetDB(),
	}
}

func (m *Match) CreateOrUpdateInBatch(ctx context.Context, mEntities []*entity.Match) ([]*entity.Match, error) {
	if len(mEntities) == 0 {
		return nil, nil
	}

	teamEntityMap := map[uint]*entity.Team{}

	var mModels []model.Match
	for _, e := range mEntities {

		teamEntityMap[e.HomeTeam.ID] = e.HomeTeam
		teamEntityMap[e.AwayTeam.ID] = e.AwayTeam

		mModels = append(mModels, model.Match{
			UID:           e.UID,
			ExternalId:    e.ExternalID,
			CompetitionID: e.CompetitionID,
			Round:         e.Round,
			HomeTeamID:    e.HomeTeam.ID,
			AwayTeamID:    e.AwayTeam.ID,
			HomeTeamScore: e.HomeTeamScore,
			AwayTeamScore: e.AwayTeamScore,
		})
	}

	err := m.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "external_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"home_team_score", "away_team_score"}),
	}).Create(&mModels).Error
	if err != nil {
		return nil, err
	}

	var result []*entity.Match
	for _, mModel := range mModels {
		result = append(result, &entity.Match{
			ID:            mModel.ID,
			UID:           mModel.UID,
			ExternalID:    mModel.ExternalId,
			CompetitionID: mModel.CompetitionID,
			Round:         mModel.Round,
			HomeTeamScore: mModel.HomeTeamScore,
			AwayTeamScore: mModel.AwayTeamScore,
			HomeTeam:      teamEntityMap[mModel.HomeTeamID],
			AwayTeam:      teamEntityMap[mModel.AwayTeamID],
		})
	}

	return result, nil
}
