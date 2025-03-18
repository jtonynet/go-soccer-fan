package dto

import "github.com/google/uuid"

type FanCreateRequest struct {
	Name     string `json:"nome" validate:"required,min=3,max=255" binding:"required" example:"Arthur Antunes Coimbra"`
	Email    string `json:"email" validate:"required,email,min=5,max=255" binding:"required" example:"zico@gmail.com"`
	TeamName string `json:"time" validate:"required,min=5,max=255" binding:"required" example:"Flamengo"`
}

type FanCreateResponse struct {
	UID      uuid.UUID `json:"id"`
	Name     string    `json:"nome"`
	Email    string    `json:"email"`
	TeamName string    `json:"time"`
	Message  string    `json:"mensagem"`
}

type FanNotification struct {
	FanUID   uuid.UUID `json:"fan_uid"`
	FanEmail string    `json:"fan_email"`
	Title    string    `json:"titulo"`
	Team     string    `json:"time"`
	Score    string    `json:"placar"`
	Message  string    `json:"mensagem"`
}
