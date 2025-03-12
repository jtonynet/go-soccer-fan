package model

import "github.com/google/uuid"

type Match struct {
	BaseModel

	UID            uuid.UUID    `gorm:"type:uuid;uniqueIndex"`
	ChampionshipID int          `gorm:"index"`
	Championship   Championship `gorm:"foreignKey:ChampionshipID"`
	HomeTeamID     int          `gorm:"index"`
	HomeTeam       Team         `gorm:"foreignKey:HomeTeamID"`
	AwayTeamID     int          `gorm:"index"`
	AwayTeam       Team         `gorm:"foreignKey:AwayTeamID"`
	Round          int
	HomeTeamScore  *int
	AwayTeamScore  *int
}
