package entity

import "github.com/google/uuid"

type Match struct {
	ID            uint
	UID           uuid.UUID
	ExternalID    string
	Round         int
	CompetitionID uint
	HomeTeam      *Team
	AwayTeam      *Team
	HomeTeamScore *int
	AwayTeamScore *int
}
