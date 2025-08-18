package review

import (
	"context"

	"github.com/google/uuid"
	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

type ReviewRepository interface {
	GetReviewToken(ctx context.Context, token string) (*model.ReviewToken, error)
	MarkReviewTokenUsed(ctx context.Context, tokenID uuid.UUID) error
	CreateReview(ctx context.Context, review model.Review) error
	HasReviewToday(ctx context.Context, userID, placeID string) (bool, error)
}
