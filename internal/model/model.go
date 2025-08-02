package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	Role         string
	Points       int
	CreatedAt    time.Time
	IsDeleted    bool
}

type Place struct {
	ID        uuid.UUID
	Name      string
	Address   string
	CreatedAt time.Time
	IsDeleted bool
}
