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

type ReviewToken struct {
	ID        uuid.UUID
	PlaceID   uuid.UUID
	Token     string
	IsUsed    bool
	ExpiresAt time.Time
}

type Review struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	PlaceID uuid.UUID
	TokenID uuid.UUID
	Content string
	Rating  int
}
