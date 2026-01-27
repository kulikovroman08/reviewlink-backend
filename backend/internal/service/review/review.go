package review

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kulikovroman08/reviewlink-backend/internal/repository"
	serviceErrors "github.com/kulikovroman08/reviewlink-backend/internal/service/errors"

	"github.com/google/uuid"
	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

const (
	RestrictionTypePointsFreeze = "review_points_freeze"
	LowRatingRestrictionReason  = "Too many negative reviews in 7 days"
	FreezeDurationDays          = 7
	LowRatingThreshold          = 3
	ReviewPeriodDays            = 7
)

type reviewService struct {
	reviewRepo      repository.ReviewRepository
	userRepo        repository.UserRepository
	placeRepo       repository.PlaceRepository
	restrictionRepo repository.UserRestrictionRepository
	replyRepo       repository.ReviewReplyRepository
}

func NewReviewService(
	reviewRepo repository.ReviewRepository,
	userRepo repository.UserRepository,
	placeRepo repository.PlaceRepository,
	restrictionRepo repository.UserRestrictionRepository,
	replyRepo repository.ReviewReplyRepository,
) *reviewService {
	return &reviewService{
		reviewRepo:      reviewRepo,
		userRepo:        userRepo,
		placeRepo:       placeRepo,
		restrictionRepo: restrictionRepo,
		replyRepo:       replyRepo,
	}
}

func (s *reviewService) SubmitReview(ctx context.Context, review model.Review, tokenStr string) error {
	tokenStr = strings.TrimSpace(tokenStr)

	if review.Rating < 1 || review.Rating > 5 {
		return serviceErrors.ErrInvalidCredentials
	}

	hasToday, err := s.reviewRepo.HasReviewToday(ctx, review.UserID.String(), review.PlaceID.String())
	if err != nil {
		return fmt.Errorf("check existing review: %w", err)
	}
	if hasToday {
		return serviceErrors.ErrTooManyReviews
	}

	var token *model.ReviewToken
	if tokenStr != "" {
		t, err := s.reviewRepo.GetReviewToken(ctx, tokenStr)
		if err != nil {
			return serviceErrors.ErrInvalidToken
		}
		if t.IsUsed {
			return serviceErrors.ErrInvalidCredentials
		}
		if t.ExpiresAt.Before(time.Now()) {
			return serviceErrors.ErrTokenExpired
		}
		if t.PlaceID != review.PlaceID {
			return serviceErrors.ErrInvalidToken
		}
		token = t
	}

	isRestricted, err := s.restrictionRepo.HasActiveRestriction(
		ctx, review.UserID.String(), RestrictionTypePointsFreeze)
	if err != nil {
		return fmt.Errorf("check restriction: %w", err)
	}

	if token != nil && !isRestricted && review.Rating == 1 {
		count, err := s.reviewRepo.CountLowRatingReviews(
			ctx, review.UserID.String(), ReviewPeriodDays)
		if err != nil {
			return fmt.Errorf("count low ratings: %w", err)
		}

		if count >= LowRatingThreshold {
			restriction := model.UserRestriction{
				ID:              uuid.New(),
				UserID:          review.UserID,
				RestrictionType: RestrictionTypePointsFreeze,
				Reason:          LowRatingRestrictionReason,
				CreatedAt:       time.Now(),
				ExpiresAt:       time.Now().AddDate(0, 0, FreezeDurationDays),
			}
			if err := s.restrictionRepo.CreateRestriction(ctx, &restriction); err != nil {
				return fmt.Errorf("create restriction: %w", err)
			}

			slog.Info("created user restriction",
				"user_id", review.UserID,
				"type", RestrictionTypePointsFreeze,
				"reason", restriction.Reason,
				"expires_at", restriction.ExpiresAt,
			)

			isRestricted = true
		}
	}

	review.ID = uuid.New()
	review.CreatedAt = time.Now().UTC()

	if token != nil {
		review.TokenID = &token.ID
	} else {
		review.TokenID = nil
	}

	if err := s.reviewRepo.CreateReview(ctx, review); err != nil {
		return fmt.Errorf("create review: %w", err)
	}

	if token == nil {
		return nil
	}

	if err := s.reviewRepo.MarkReviewTokenUsed(ctx, token.ID.String()); err != nil {
		return fmt.Errorf("mark token used: %w", err)
	}

	if isRestricted {
		fmt.Printf("User %s has active restriction: skip points\n", review.UserID)
		return nil
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

	repliesIDs := make([]string, 0, len(reviews))
	for _, r := range reviews {
		repliesIDs = append(repliesIDs, r.ID.String())
	}

	replies, err := s.replyRepo.GetByReviewIDs(ctx, repliesIDs)
	if err != nil {
		return nil, fmt.Errorf("get replies: %w", err)
	}

	for i := range reviews {
		if rr := replies[reviews[i].ID.String()]; rr != nil {
			reviews[i].Reply = rr
		}
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

func (s *reviewService) DeleteReview(ctx context.Context, reviewID, userID string) error {
	err := s.reviewRepo.DeleteReview(ctx, reviewID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return serviceErrors.ErrReviewNotFound
		}
		return fmt.Errorf("delete review: %w", err)
	}
	return nil
}
