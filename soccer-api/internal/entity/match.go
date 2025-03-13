package entity

import "github.com/google/uuid"

type Match struct {
	ID            uint
	UID           uuid.UUID
	Round         int
	CompetitionID int
	HomeTeam      *Team
	AwayTeam      *Team
	HomeTeamScore *int
	AwayTeamScore *int
}
