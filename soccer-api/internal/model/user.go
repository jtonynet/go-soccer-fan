package model

import "github.com/google/uuid"

type User struct {
	BaseModel

	UID      uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	Username string    `gorm:"size:255;not null;unique" json:"username"`
	Password string    `gorm:"size:255;not null;" json:"password"`
	Name     string    `gorm:"size:255;not null;" json:"name"`
	Email    string    `gorm:"size:255;not null;" json:"email"`
}
