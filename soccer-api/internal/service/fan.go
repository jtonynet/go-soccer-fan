package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/dto"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/repository"
)

type Fan struct {
	fRepo repository.Fan
}

func NewFan(fRepo repository.Fan) *Fan {
	return &Fan{fRepo}
}

func (f *Fan) Create(fReq *dto.FanCreateRequest) (*dto.FanCreateResponse, error) {
	fEntity := mapFanCreateRequestDTOtoEntity(fReq)

	createdEntity, err := f.fRepo.Create(context.Background(), fEntity)
	if err != nil {
		return nil, err
	}

	return mapFanEntityToCreateResponseDTO(createdEntity), nil
}

func mapFanCreateRequestDTOtoEntity(fReq *dto.FanCreateRequest) *entity.Fan {
	return &entity.Fan{
		UID:   uuid.New(),
		Name:  fReq.Name,
		Email: fReq.Email,
		Team: &entity.Team{
			Name: fReq.TeamName,
		},
	}
}

func mapFanEntityToCreateResponseDTO(fEntity *entity.Fan) *dto.FanCreateResponse {
	return &dto.FanCreateResponse{
		UID:      fEntity.UID,
		Name:     fEntity.Name,
		Email:    fEntity.Email,
		TeamName: fEntity.Team.Name,
		Message:  "Cadastro realizado com sucesso",
	}
}
