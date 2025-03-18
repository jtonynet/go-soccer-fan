package dto

import "github.com/google/uuid"

type UserCreateRequest struct {
	UserName string `json:"usuario" validate:"required,min=3,max=255" binding:"required" example:"admin"`
	Password string `json:"senha" validate:"required,min=5,max=255" binding:"required" example:"admin"`
	Name     string `json:"nome" validate:"required,min=3,max=255" binding:"required" example:"Edson Arantes do Nascimento"`
	Email    string `json:"email" validate:"required,email,min=5,max=255" binding:"required" example:"pele@soccerfan.com"`
}

type UserCreateResponse struct {
	UID      uuid.UUID `json:"id"`
	UserName string    `json:"usuario"`
	Name     string    `json:"nome"`
	Email    string    `json:"email"`
}

type UserLoginRequest struct {
	UserName string `json:"usuario" validate:"required,min=3,max=255" binding:"required" example:"admin"`
	Password string `json:"senha" validate:"required,min=5,max=255" binding:"required" example:"admin"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}
