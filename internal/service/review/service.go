package review

import (
	"context"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

type ReviewService interface {
	SubmitReview(ctx context.Context, review model.Review, token string) error
}
