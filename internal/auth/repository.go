package auth

import (
	"context"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, users *User) error
}
