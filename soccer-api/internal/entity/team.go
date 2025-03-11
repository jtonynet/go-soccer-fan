package entity

import "github.com/google/uuid"

type Team struct {
	ID   int
	UID  uuid.UUID
	Name string
}
