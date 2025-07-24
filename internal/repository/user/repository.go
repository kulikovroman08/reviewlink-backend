package user

import (
	"context"

	"github.com/kulikovroman08/reviewlink-backend/internal/service/user"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	FindByID(ctx context.Context, id string) (*user.User, error)
	CreateUser(ctx context.Context, users *user.User) error
}
