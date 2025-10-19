package repository

import (
	"context"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

//go:generate go run go.uber.org/mock/mockgen -source=repository.go -destination=../tests/integration/mocks/repository_mocks.go -package=mocks

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, userID string) (*model.User, error)
	FindAnyByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, users *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	SoftDeleteUser(ctx context.Context, userID string) error
	AddPoints(ctx context.Context, userID string, points int) error
}

type PlaceRepository interface {
	CreatePlace(ctx context.Context, place *model.Place) error
	GetByID(ctx context.Context, placeID string) (*model.Place, error)
}

type ReviewRepository interface {
	GetReviewToken(ctx context.Context, token string) (*model.ReviewToken, error)
	MarkReviewTokenUsed(ctx context.Context, tokenID string) error
	CreateReview(ctx context.Context, review model.Review) error
	HasReviewToday(ctx context.Context, userID, placeID string) (bool, error)
	FindReviews(ctx context.Context, placeID string, filter model.ReviewFilter) ([]model.Review, error)
	UpdateReview(ctx context.Context, reviewID, userID string, content string, rating int) error
	DeleteReview(ctx context.Context, reviewID, userID string) error
}

type TokenRepository interface {
	CreateTokens(ctx context.Context, tokens []model.ReviewToken) error
	CountActiveTokens(ctx context.Context, placeID string) (int, error)
}

type AdminRepository interface {
	GetAdminStats(ctx context.Context) (*model.AdminStats, error)
}
