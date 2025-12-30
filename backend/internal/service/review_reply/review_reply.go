package review_reply

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
	"github.com/kulikovroman08/reviewlink-backend/internal/repository"
	serviceErrors "github.com/kulikovroman08/reviewlink-backend/internal/service/errors"
)

type service struct {
	reviewRepo repository.ReviewRepository
	placeRepo  repository.PlaceRepository
	replyRepo  repository.ReviewReplyRepository
}

func NewReviewReplyService(
	reviewRepo repository.ReviewRepository,
	placeRepo repository.PlaceRepository,
	replyRepo repository.ReviewReplyRepository,
) *service {
	return &service{
		reviewRepo: reviewRepo,
		placeRepo:  placeRepo,
		replyRepo:  replyRepo,
	}
}

func (s *service) ReplyToReview(ctx context.Context, adminID string, reviewID string, content string) (*model.ReviewReply, error) {
	if content == "" {
		return nil, serviceErrors.ErrInvalidCredentials
	}

	rev, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, serviceErrors.ErrReviewNotFound
		}
		return nil, fmt.Errorf("get review: %w", err)
	}

	isOwner, err := s.placeRepo.IsOwner(ctx, rev.PlaceID.String(), adminID)
	if err != nil {
		return nil, fmt.Errorf("check owner: %w", err)
	}
	if !isOwner {
		return nil, serviceErrors.ErrAccessDenied // добавь если нет
	}

	if _, err := s.replyRepo.GetByReviewID(ctx, reviewID); err == nil {
		return nil, serviceErrors.ErrReplyAlreadyExists
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("check existing reply: %w", err)
	}

	reply := &model.ReviewReply{
		ID:        uuid.New(),
		ReviewID:  rev.ID,
		AdminID:   uuid.MustParse(adminID),
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := s.replyRepo.CreateReply(ctx, reply); err != nil {
		return nil, fmt.Errorf("create reply: %w", err)
	}

	return reply, nil
}

func (s *service) UpdateReply(ctx context.Context, adminID string, reviewID string, content string) (*model.ReviewReply, error) {
	if content == "" {
		return nil, serviceErrors.ErrInvalidCredentials
	}

	rev, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, serviceErrors.ErrReviewNotFound
		}
		return nil, fmt.Errorf("get review: %w", err)
	}

	isOwner, err := s.placeRepo.IsOwner(ctx, rev.PlaceID.String(), adminID)
	if err != nil {
		return nil, fmt.Errorf("check owner: %w", err)
	}
	if !isOwner {
		return nil, serviceErrors.ErrAccessDenied
	}

	updated, err := s.replyRepo.UpdateReply(ctx, reviewID, content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, serviceErrors.ErrReplyNotFound
		}
		return nil, fmt.Errorf("update reply: %w", err)
	}

	return updated, nil
}
