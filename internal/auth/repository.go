package auth

import (
	"context"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	CreateUsers(ctx context.Context, users *User) error
}
