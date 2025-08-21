package review

import (
	"context"
	"fmt"
	"time"

	"github.com/kulikovroman08/reviewlink-backend/internal/repository/place"
	"github.com/kulikovroman08/reviewlink-backend/internal/repository/review"
	"github.com/kulikovroman08/reviewlink-backend/internal/repository/user"

	"github.com/google/uuid"
	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

type reviewService struct {
	reviewRepo review.ReviewRepository
	userRepo   user.UserRepository
	placeRepo  place.PlaceRepository
}

func NewReviewService(
	reviewRepo review.ReviewRepository,
	userRepo user.UserRepository,
	placeRepo place.PlaceRepository,
) ReviewService {
	return &reviewService{
		reviewRepo: reviewRepo,
		userRepo:   userRepo,
		placeRepo:  placeRepo,
	}
}

func (s *reviewService) SubmitReview(ctx context.Context, review model.Review, tokenStr string) error {
	if tokenStr == "" {
		return fmt.Errorf("token required")
	}

	if review.Rating < 1 || review.Rating > 5 {
		return fmt.Errorf("invalid rating")
	}

	token, err := s.reviewRepo.GetReviewToken(ctx, tokenStr)
	if err != nil {
		return fmt.Errorf("get token: %w", err)
	}
	if token.IsUsed {
		return fmt.Errorf("token is used")
	}

	if token.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("token is expired")
	}

	hashToday, err := s.reviewRepo.HasReviewToday(ctx, review.UserID.String(), token.PlaceID.String())
	if err != nil {
		return fmt.Errorf("check existing review: %w", err)
	}
	if hashToday {
		return fmt.Errorf("review already submitted today for this place")
	}

	review.ID = uuid.New()
	review.TokenID = token.ID

	if err := s.reviewRepo.CreateReview(ctx, review); err != nil {
		return fmt.Errorf("create review: %w", err)
	}

	if err := s.reviewRepo.MarkReviewTokenUsed(ctx, token.ID.String()); err != nil {
		return fmt.Errorf("mark token used: %w", err)
	}

	var points int
	switch review.Rating {
	case 5:
		points = 10
	case 4:
		points = 5
	default:
		points = 0
	}

	if points > 0 {
		if err := s.userRepo.AddPoints(ctx, review.UserID.String(), points); err != nil {
			return fmt.Errorf("add points: %w", err)
		}
	}

	return nil
}
