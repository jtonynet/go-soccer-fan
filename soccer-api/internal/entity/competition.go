package entity

import "github.com/google/uuid"

type Competition struct {
	ID         uint
	UID        uuid.UUID
	ExternalId string
	Name       string
	Season     string
}
