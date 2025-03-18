package dto

import "github.com/google/uuid"

type UserCreateRequest struct {
	UserName string `json:"usuario" binding:"required"`
	Password string `json:"senha" binding:"required"`
	Name     string `json:"nome" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserCreateResponse struct {
	UID      uuid.UUID `json:"id"`
	UserName string    `json:"usuario"`
	Name     string    `json:"nome"`
	Email    string    `json:"email"`
}

type UserLoginRequest struct {
	UserName string `json:"usuario" binding:"required"`
	Password string `json:"senha" binding:"required"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}
