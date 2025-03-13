package model

import "github.com/google/uuid"

type Team struct {
	BaseModel

	UID        uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	ExternalId string    `gorm:"type:varchar(255);uniqueIndex"`
	Name       string    `gorm:"type:varchar(255);index"`
	FullName   string    `gorm:"type:varchar(255)"`
}
