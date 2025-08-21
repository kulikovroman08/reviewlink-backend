package user

import (
	"context"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, userID string) (*model.User, error)
	FindAnyByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, users *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	SoftDeleteUser(ctx context.Context, userID string) error
	AddPoints(ctx context.Context, userID string, points int) error
}
