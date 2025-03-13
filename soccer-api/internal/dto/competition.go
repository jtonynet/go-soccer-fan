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

type CompetitionResponseExternalAPI struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type CompetitionResponseListExternalAPI struct {
	Count        int                              `json:"count"`
	Filters      map[string]interface{}           `json:"filters"`
	Competitions []CompetitionResponseExternalAPI `json:"competitions"`
}
