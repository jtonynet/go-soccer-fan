package model

import "github.com/google/uuid"

type Match struct {
	BaseModel

	UID           uuid.UUID   `gorm:"type:uuid;uniqueIndex"`
	ExternalId    string      `gorm:"type:varchar(255);uniqueIndex"`
	CompetitionID int         `gorm:"index"`
	Competition   Competition `gorm:"foreignKey:CompetitionID"`
	HomeTeamID    int         `gorm:"index"`
	HomeTeam      Team        `gorm:"foreignKey:HomeTeamID"`
	AwayTeamID    int         `gorm:"index"`
	AwayTeam      Team        `gorm:"foreignKey:AwayTeamID"`
	Round         int
	HomeTeamScore *int
	AwayTeamScore *int
}
