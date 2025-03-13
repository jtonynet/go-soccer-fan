package entity

import "github.com/google/uuid"

type Team struct {
	ID         uint
	UID        uuid.UUID
	ExternalId string
	Name       string
	FullName   string
}
