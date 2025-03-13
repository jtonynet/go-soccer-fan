package entity

import "github.com/google/uuid"

type Team struct {
	ID   uint
	UID  uuid.UUID
	Name string
}
