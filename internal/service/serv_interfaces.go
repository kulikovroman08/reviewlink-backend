package service

import (
	"context"

	"github.com/kulikovroman08/reviewlink-backend/internal/service/user/model"
)

type UserService interface {
	Signup(ctx context.Context, name, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	GetProfile(ctx context.Context, userID string) (*model.User, error)
}
