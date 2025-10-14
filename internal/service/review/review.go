package review

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
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

func (s *reviewService) GetReviews(ctx context.Context, placeID string, filter model.ReviewFilter) ([]model.Review, error) {
	_, err := s.placeRepo.GetByID(ctx, placeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, serviceErrors.ErrPlaceNotFound
		}
		return nil, fmt.Errorf("check place existence: %w", err)
	}

	reviews, err := s.reviewRepo.FindReviews(ctx, placeID, filter)
	if err != nil {
		return nil, fmt.Errorf("find reviews: %w", err)
	}

	return reviews, nil
}

func (s *reviewService) UpdateReview(ctx context.Context, reviewID, userID string, content string, rating int) error {
	if rating < 1 || rating > 5 {
		return serviceErrors.ErrInvalidRating
	}

	err := s.reviewRepo.UpdateReview(ctx, reviewID, userID, content, rating)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return serviceErrors.ErrReviewNotFound
		}
		return err
	}

	return nil
}
