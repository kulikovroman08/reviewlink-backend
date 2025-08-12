package review

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
	"github.com/kulikovroman08/reviewlink-backend/internal/model"
	usersvc "github.com/kulikovroman08/reviewlink-backend/internal/service/user"
)

type Service struct {
	repo     ReviewRepository
	userRepo usersvc.UserRepository
}

func NewService(repo ReviewRepository, userRepo usersvc.UserRepository) *Service {
	return &Service{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *Service) SubmitReview(ctx context.Context, userID string, req dto.SubmitReviewRequest) error {
	if req.Token == "" {
		return fmt.Errorf("token required")
	}

	if req.Rating < 1 || req.Rating > 5 {
		return fmt.Errorf("invalid rating")
	}

	token, err := s.repo.GetReviewToken(ctx, req.Token)
	if err != nil {
		fmt.Println("ERR GET TOKEN:", err)
		return fmt.Errorf("get token: %w", err)
	}
	fmt.Println("TOKEN FOUND:", token)
	fmt.Println("IS USED?", token.IsUsed)
	fmt.Println("EXPIRES AT:", token.ExpiresAt)
	fmt.Println("NOW:", time.Now())
	if token.IsUsed {
		return fmt.Errorf("token is used")
	}

	if token.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("token is expired")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	hashToday, err := s.repo.HasReviewToday(ctx, userUUID, token.PlaceID)
	if err != nil {
		return fmt.Errorf("check existing review: %w", err)
	}
	if hashToday {
		return fmt.Errorf("review already submitted today for this place")
	}

	reviewID := uuid.New()
	review := model.Review{
		ID:      reviewID,
		UserID:  userUUID,
		PlaceID: token.PlaceID,
		TokenID: token.ID,
		Rating:  req.Rating,
		Content: req.Content,
	}

	if err := s.repo.CreateReview(ctx, review); err != nil {
		return fmt.Errorf("create review: %w", err)
	}

	if err := s.repo.MarkReviewTokenUsed(ctx, token.ID, time.Now()); err != nil {
		return fmt.Errorf("mark token used: %w", err)
	}

	var points int
	switch req.Rating {
	case 5:
		points = 10
	case 4:
		points = 5
	default:
		points = 0
	}

	if points > 0 {
		if err := s.userRepo.AddPoints(ctx, userUUID, points); err != nil {
			return fmt.Errorf("add points: %w", err)
		}
	}

	return nil
}
