package entity

import "github.com/google/uuid"

type Match struct {
	ID            int
	UID           uuid.UUID
	Round         int
	CompetitionID int
	HomeTeam      *Team
	AwayTeam      *Team
	HomeTeamScore *int
	AwayTeamScore *int
}
