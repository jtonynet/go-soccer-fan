package dto

import "github.com/google/uuid"

type UserCreateRequest struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserCreateResponse struct {
	UID      uuid.UUID
	UserName string
	Name     string
	Email    string
}
