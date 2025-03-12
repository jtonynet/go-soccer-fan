package model

import "github.com/google/uuid"

type Competition struct {
	BaseModel

	UID    uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	Name   string    `gorm:"type:varchar(255)"`
	Season string    `gorm:"type:varchar(255)"`
}
