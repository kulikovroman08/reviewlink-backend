package auth

import "time"

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
