package user

import (
	"context"

	"github.com/kulikovroman08/reviewlink-backend/internal/service/user/model"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	CreateUser(ctx context.Context, users *model.User) error
}
