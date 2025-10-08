package service

import (
	"context"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

//go:generate go run go.uber.org/mock/mockgen -source=service.go -destination=../tests/integration/mocks/service_mocks.go -package=mocks

type UserService interface {
	Signup(ctx context.Context, name, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	GetUser(ctx context.Context, userID string) (*model.User, error)
	UpdateUser(ctx context.Context, user model.User, password string) (*model.User, error)
	DeleteUser(ctx context.Context, userID string) error
}

type PlaceService interface {
	CreatePlace(ctx context.Context, place model.Place) (*model.Place, error)
}

type ReviewService interface {
	SubmitReview(ctx context.Context, review model.Review, token string) error
	GetReviews(ctx context.Context, placeID string, filter model.ReviewFilter) ([]model.Review, error)
}

type TokenService interface {
	GenerateTokens(ctx context.Context, placeID string, count int) (*model.GenerateTokensResult, error)
}
