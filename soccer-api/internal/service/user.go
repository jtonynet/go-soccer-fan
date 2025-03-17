package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/dto"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/repository"
)

type User struct {
	uRepo repository.User
}

func NewUser(uRepo repository.User) *User {
	return &User{uRepo}
}

func (u *User) Create(uReq *dto.UserCreateRequest) (*dto.UserCreateResponse, error) {
	uEntity := mapUserCreateRequestDTOtoEntity(uReq)

	createdEntity, err := u.uRepo.Create(context.Background(), uEntity)
	if err != nil {
		return nil, err
	}

	return mapUserEntityToCreateResponseDTO(createdEntity), nil
}

func mapUserCreateRequestDTOtoEntity(uReq *dto.UserCreateRequest) *entity.User {
	return &entity.User{
		UID:      uuid.New(),
		UserName: uReq.UserName,
		Password: uReq.Password,
		Name:     uReq.Name,
		Email:    uReq.Email,
	}
}

func mapUserEntityToCreateResponseDTO(uEntity *entity.User) *dto.UserCreateResponse {
	return &dto.UserCreateResponse{
		UID:      uEntity.UID,
		UserName: uEntity.UserName,
		Name:     uEntity.Name,
		Email:    uEntity.Email,
	}
}
