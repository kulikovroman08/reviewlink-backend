package review

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

type Repository interface {
	GetReviewToken(ctx context.Context, token string) (*model.ReviewToken, error)
	MarkReviewTokenUsed(ctx context.Context, tokenID uuid.UUID, usedAt time.Time) error
	CreateReview(ctx context.Context, review model.Review) error
}
