package review

import (
	"context"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

type ReviewRepository interface {
	GetReviewToken(ctx context.Context, token string) (*model.ReviewToken, error)
	MarkReviewTokenUsed(ctx context.Context, tokenID string) error
	CreateReview(ctx context.Context, review model.Review) error
	HasReviewToday(ctx context.Context, userID, placeID string) (bool, error)
}
