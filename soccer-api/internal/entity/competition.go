package entity

import "github.com/google/uuid"

type Competition struct {
	ID     uint
	UID    uuid.UUID
	Name   string
	Season string
}
