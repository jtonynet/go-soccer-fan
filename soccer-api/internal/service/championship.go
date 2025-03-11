package service

import (
	"context"

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

	cEntityList, err := c.cRepo.FindAll(context.Background())
	if err != nil {
		return nil, err
	}

	result := dto.ChampionshipResponseList{}
	for _, cEntity := range cEntityList {
		result.Championships = append(
			result.Championships,
			mapChampionshipEntityToResponseDTO(cEntity),
		)
	}

	return &result, nil
}

func mapChampionshipEntityToResponseDTO(ce *entity.Championship) *dto.ChampionshipResponse {
	return &dto.ChampionshipResponse{
		UID:    ce.UID,
		Name:   ce.Name,
		Season: ce.Season,
	}
}
