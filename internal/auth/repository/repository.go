package repository

import (
	"context"

	"github.com/kulikovroman08/reviewlink-backend/internal/auth"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*auth.User, error)
	FindByID(ctx context.Context, id string) (*auth.User, error)
	CreateUser(ctx context.Context, users *auth.User) error
}
