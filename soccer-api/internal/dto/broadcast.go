package dto

type BroadcastSendRequest struct {
	Type     string `json:"tipo" validate:"required,oneof=inicio fim" binding:"required" example:"fim"`
	TeamName string `json:"time" validate:"required,min=5,max=255" binding:"required" example:"Flamengo"`
	Score    string `json:"placar" validate:"omitempty"`
	Message  string `json:"mensagem" validate:"required,min=5,max=255" binding:"required" example:"Flamengo"`
}

type BroadcastResponse struct {
	Message string `json:"mensagem"`
}
