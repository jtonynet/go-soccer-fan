package model

import "github.com/google/uuid"

type Team struct {
	BaseModel

	UID  uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	Name string    `gorm:"type:varchar(255);index"`
}
