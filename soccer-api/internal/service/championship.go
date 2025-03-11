package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/dto"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/repository"
)

type Championship struct {
	cRepo repository.Championship
}

func NewChampionship(cRepo repository.Championship) *Championship {
	return &Championship{cRepo}
}

func (c *Championship) FindAll() (*dto.ChampionshipResponseList, error) {

	cEntities, err := c.cRepo.FindAll(context.Background())
	if err != nil {
		return nil, err
	}

	result := mapChampionshipEntitiesToResponseListDTO(cEntities)
	return result, nil
}

func (c *Championship) FindMatchsByChampionshipUID(uid uuid.UUID) (*dto.RoundResponseList, error) {
	matchEntities, err := c.cRepo.FindMatchsByChampionshipUID(context.Background(), uid)
	if err != nil {
		return nil, err
	}

	return mapMatchEntitiesToResponseListDTO(matchEntities), nil
}

func mapMatchEntitiesToResponseListDTO(mEntities []*entity.Match) *dto.RoundResponseList {
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

	return &dto.RoundResponseList{
		Rounds: rounds,
	}
}

func mapChampionshipEntitiesToResponseListDTO(cEntities []*entity.Championship) *dto.ChampionshipResponseList {
	result := dto.ChampionshipResponseList{}
	for _, cEntity := range cEntities {
		result.Championships = append(
			result.Championships,
			&dto.ChampionshipResponse{
				UID:    cEntity.UID,
				Name:   cEntity.Name,
				Season: cEntity.Season,
			},
		)
	}

	return &result
}
