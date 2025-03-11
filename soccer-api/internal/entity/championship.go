package entity

import "github.com/google/uuid"

type Championship struct {
	ID     int
	UID    uuid.UUID
	Name   string
	Season string
}
