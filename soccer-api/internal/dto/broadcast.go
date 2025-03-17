package dto

const (
	TypeStart = "inicio"
	TypeEnd   = "fim"
)

type BroadcastSendRequest struct {
	Type    string `json:"tipo"`
	Team    string `json:"time"`
	Score   string `json:"placar"`
	Message string `json:"mensagem"`
}

type BroadcastResponse struct {
	Message string `json:"mensagem"`
}
