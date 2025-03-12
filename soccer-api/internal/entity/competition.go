package entity

import "github.com/google/uuid"

type Competition struct {
	ID     int
	UID    uuid.UUID
	Name   string
	Season string
}
