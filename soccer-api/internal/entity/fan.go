package entity

import "github.com/google/uuid"

type Fan struct {
	ID    uint
	UID   uuid.UUID
	Name  string
	Email string
	Team  *Team
}
