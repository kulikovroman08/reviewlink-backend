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
	GetUserStats(ctx context.Context, userID string) (*model.UserStats, error)
}

type PlaceService interface {
	CreatePlace(ctx context.Context, place model.Place) (*model.Place, error)
	GetAllPlaces(ctx context.Context) ([]model.Place, error)
}

type ReviewService interface {
	SubmitReview(ctx context.Context, review model.Review, token string) error
	GetReviews(ctx context.Context, placeID string, filter model.ReviewFilter) ([]model.Review, error)
	UpdateReview(ctx context.Context, reviewID, userID string, content string, rating int) error
	DeleteReview(ctx context.Context, reviewID, userID string) error
}

type TokenService interface {
	GenerateTokens(ctx context.Context, placeID string, count int) (*model.GenerateTokensResult, error)
	CheckAndRefillTokens(ctx context.Context, placeID string) error
}

type AdminService interface {
	GetStats(ctx context.Context) (*model.AdminStats, error)
}

type LeaderboardService interface {
	GetUserLeaderboard(ctx context.Context, limit int, filter model.LeaderboardFilter) ([]model.LeaderboardEntry, error)
	GetPlaceLeaderboard(ctx context.Context, limit int, filter model.LeaderboardFilter) ([]model.LeaderboardEntry, error)
	GetBonusLeaderboard(ctx context.Context) ([]model.BonusLeaderboardEntry, error)
}

type BonusService interface {
	RedeemBonus(ctx context.Context, userID, rewardType string) (*model.BonusReward, error)
	GetUserBonuses(ctx context.Context, userID string) ([]model.BonusReward, error)
	ValidateBonus(ctx context.Context, qrToken string) error
}
