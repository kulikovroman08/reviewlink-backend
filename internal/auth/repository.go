package auth

import (
	"context"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	CreateUser(ctx context.Context, users *User) error
}
