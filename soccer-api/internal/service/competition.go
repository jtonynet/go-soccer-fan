package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/dto"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/repository"
)

type Competition struct {
	cRepo repository.Competition
}

func NewCompetition(cRepo repository.Competition) *Competition {
	return &Competition{cRepo}
}

func (c *Competition) FindAll() (*dto.CompetitionResponseList, error) {

	cEntities, err := c.cRepo.FindAll(context.Background())
	if err != nil {
		return nil, err
	}

	result := mapCompetitionEntitiesToResponseListDTO(cEntities)
	return result, nil
}

func (c *Competition) FindMatchsByCompetitionUID(uid uuid.UUID) (*dto.MatchResponseList, error) {
	matchEntities, err := c.cRepo.FindMatchsByCompetitionUID(context.Background(), uid)
	if err != nil {
		return nil, err
	}

	return mapMatchEntitiesToResponseListDTO(matchEntities), nil
}

func mapMatchEntitiesToResponseListDTO(mEntities []*entity.Match) *dto.MatchResponseList {
	roundsMap := map[int]dto.MatchResponseListPerRound{}
	for _, mEntity := range mEntities {
		var matchsPerRound dto.MatchResponseListPerRound
		var exists bool

		if matchsPerRound, exists = roundsMap[mEntity.Round]; !exists {
			matchsPerRound = dto.MatchResponseListPerRound{
				Round: mEntity.Round,
			}
		}

		score := "-"
		if mEntity.HomeTeamScore != nil && mEntity.AwayTeamScore != nil {
			score = fmt.Sprintf(`%v-%v`, *mEntity.HomeTeamScore, *mEntity.AwayTeamScore)
		}

		matchsPerRound.Matchs = append(matchsPerRound.Matchs, &dto.MatchResponse{
			HomeTeamName: mEntity.HomeTeam.Name,
			AwayTeamName: mEntity.AwayTeam.Name,
			Score:        score,
		})

		roundsMap[mEntity.Round] = matchsPerRound
	}

	var rounds []*dto.MatchResponseListPerRound
	for _, matchsPerRound := range roundsMap {
		rounds = append(rounds, &matchsPerRound)
	}

	if len(rounds) == 0 {
		rounds = []*dto.MatchResponseListPerRound{}
	}

	return &dto.MatchResponseList{
		Rounds: rounds,
	}
}

func mapCompetitionEntitiesToResponseListDTO(cEntities []*entity.Competition) *dto.CompetitionResponseList {
	result := dto.CompetitionResponseList{}
	for _, cEntity := range cEntities {
		result.Competitions = append(
			result.Competitions,
			&dto.CompetitionResponse{
				UID:    cEntity.UID,
				Name:   cEntity.Name,
				Season: cEntity.Season,
			},
		)
	}

	if len(result.Competitions) == 0 {
		result.Competitions = []*dto.CompetitionResponse{}
	}

	return &result
}
