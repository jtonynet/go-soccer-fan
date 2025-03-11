package dto

import "github.com/google/uuid"

type ChampionshipResponse struct {
	UID    uuid.UUID `json:"id"`
	Name   string    `json:"nome"`
	Season string    `json:"temporada"`
}

type ChampionshipResponseList struct {
	Championships []*ChampionshipResponse `json:"campeonatos,omitempty"`
}
