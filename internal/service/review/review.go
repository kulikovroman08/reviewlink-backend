package review

import (
	"context"
	"fmt"
	"time"

	serviceErrors "github.com/kulikovroman08/reviewlink-backend/internal/service/errors"

	"github.com/kulikovroman08/reviewlink-backend/internal/repository"
	"github.com/kulikovroman08/reviewlink-backend/internal/service"

	"github.com/google/uuid"
	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

type reviewService struct {
	reviewRepo repository.ReviewRepository
	userRepo   repository.UserRepository
	placeRepo  repository.PlaceRepository
}

func NewReviewService(
	reviewRepo repository.ReviewRepository,
	userRepo repository.UserRepository,
	placeRepo repository.PlaceRepository,
) service.ReviewService {
	return &reviewService{
		reviewRepo: reviewRepo,
		userRepo:   userRepo,
		placeRepo:  placeRepo,
	}
}

func (s *reviewService) SubmitReview(ctx context.Context, review model.Review, tokenStr string) error {
	if tokenStr == "" {
		return serviceErrors.ErrInvalidCredentials
	}

	if review.Rating < 1 || review.Rating > 5 {
		return serviceErrors.ErrInvalidCredentials
	}

	token, err := s.reviewRepo.GetReviewToken(ctx, tokenStr)
	if err != nil {
		return fmt.Errorf("get token: %w", err)
	}
	if token.IsUsed {
		return serviceErrors.ErrInvalidCredentials
	}

	if token.ExpiresAt.Before(time.Now()) {
		return serviceErrors.ErrTokenExpired
	}

	hashToday, err := s.reviewRepo.HasReviewToday(ctx, review.UserID.String(), token.PlaceID.String())
	if err != nil {
		return fmt.Errorf("check existing review: %w", err)
	}
	if hashToday {
		return serviceErrors.ErrInvalidCredentials
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
