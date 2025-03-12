package dto

import "github.com/google/uuid"

type CompetitionResponse struct {
	UID    uuid.UUID `json:"id"`
	Name   string    `json:"nome"`
	Season string    `json:"temporada"`
}

type CompetitionResponseList struct {
	Competitions []*CompetitionResponse `json:"campeonatos"`
}
