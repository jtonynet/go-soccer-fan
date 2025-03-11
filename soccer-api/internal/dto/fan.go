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
