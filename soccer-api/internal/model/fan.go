package model

import "github.com/google/uuid"

type Fan struct {
	BaseModel

	UID    uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	Name   string    `gorm:"type:varchar(255)"`
	Email  string    `gorm:"type:varchar(255)"`
	TeamID int       `gorm:"index"`
	Team   Team      `gorm:"foreignKey:TeamID"`
}
