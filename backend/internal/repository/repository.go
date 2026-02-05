package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

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
	RedeemPoints(ctx context.Context, userID string, points int) error
	SetRole(ctx context.Context, userID string, role string) error
}

type PlaceRepository interface {
	CreatePlace(ctx context.Context, place *model.Place) error
	GetByID(ctx context.Context, placeID string) (*model.Place, error)
	GetPlacesByOwner(ctx context.Context, ownerID string) ([]model.Place, error)
	IsOwner(ctx context.Context, placeID string, ownerID string) (bool, error)
	GetPublicMetaBySourceIDs(ctx context.Context, source string, sourceIDs []string) (map[string]model.PublicPlaceMeta, error)
	EnsureFromPublic(ctx context.Context, source, sourceID, name string) (string, error)
	SetOwner(ctx context.Context, placeID string, ownerID string) error
}

type ReviewRepository interface {
	GetReviewToken(ctx context.Context, token string) (*model.ReviewToken, error)
	MarkReviewTokenUsed(ctx context.Context, tokenID string) error
	CreateReview(ctx context.Context, review model.Review) error
	HasReviewToday(ctx context.Context, userID, placeID string) (bool, error)
	FindReviews(ctx context.Context, placeID string, filter model.ReviewFilter) ([]model.Review, error)
	UpdateReview(ctx context.Context, reviewID, userID string, content string, rating int) error
	DeleteReview(ctx context.Context, reviewID, userID string) error
	CountLowRatingReviews(ctx context.Context, userID string, days int) (int, error)
	CountUserReviews(ctx context.Context, userID string) (int, error)
	AvgUserRating(ctx context.Context, userID string) (float64, error)
	GetByID(ctx context.Context, reviewID string) (*model.Review, error)
}

type ReviewReplyRepository interface {
	CreateReply(ctx context.Context, reply *model.ReviewReply) error
	UpdateReply(ctx context.Context, reviewID string, content string) (*model.ReviewReply, error)
	GetByReviewID(ctx context.Context, reviewID string) (*model.ReviewReply, error)
	GetByReviewIDs(ctx context.Context, reviewIDs []string) (map[string]*model.ReviewReply, error)
}

type TokenRepository interface {
	CreateTokens(ctx context.Context, tokens []model.ReviewToken) error
	CountActiveTokens(ctx context.Context, placeID string) (int, error)
}

type AdminRepository interface {
	GetAdminStats(ctx context.Context) (*model.AdminStats, error)
}

type LeaderboardRepository interface {
	GetTopUsers(ctx context.Context, limit int, filter model.LeaderboardFilter) ([]model.LeaderboardEntry, error)
	GetTopPlaces(ctx context.Context, limit int, filter model.LeaderboardFilter) ([]model.LeaderboardEntry, error)
	GetTopBonusUsers(ctx context.Context) ([]model.BonusLeaderboardEntry, error)
}

type BonusRepository interface {
	CreateBonus(ctx context.Context, bonus *model.BonusReward) error
	GetBonusesByUser(ctx context.Context, userID string) ([]model.BonusReward, error)
	MarkBonusUsed(ctx context.Context, qrToken string) error
	GetByQRToken(ctx context.Context, qrToken string) (*model.BonusReward, error)
}

type UserRestrictionRepository interface {
	HasActiveRestriction(ctx context.Context, userID, restrictionType string) (bool, error)
	CreateRestriction(ctx context.Context, restriction *model.UserRestriction) error
}

type ReviewVoteRepository interface {
	GetVoteForUpdate(ctx context.Context, tx pgx.Tx, reviewID, userID string) (*int16, error)
	UpsertVote(ctx context.Context, tx pgx.Tx, reviewID, userID string, value int16) error
	DeleteVote(ctx context.Context, tx pgx.Tx, reviewID, userID string) (bool, error)
	UpdateCounters(ctx context.Context, tx pgx.Tx, reviewID string, deltaHelpful int32, deltaUnhelpful int32) error
}

type OwnerRequestRepository interface {
	Create(ctx context.Context, r *model.OwnerRequest) error
	GetByID(ctx context.Context, id string) (*model.OwnerRequest, error)
	ListPending(ctx context.Context, limit int) ([]model.OwnerRequest, error)
	Approve(ctx context.Context, id string, reviewedBy string, comment *string) error
	Reject(ctx context.Context, id string, reviewedBy string, comment *string) error
}
