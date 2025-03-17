package dto

import "github.com/google/uuid"

type FanCreateRequest struct {
	Name     string `json:"nome" `
	Email    string `json:"email"`
	TeamName string `json:"time" `
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
