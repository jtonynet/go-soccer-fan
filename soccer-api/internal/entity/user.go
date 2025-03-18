package entity

import "github.com/google/uuid"

type User struct {
	ID       uint
	UID      uuid.UUID
	UserName string
	Password string
	Name     string
	Email    string
}
