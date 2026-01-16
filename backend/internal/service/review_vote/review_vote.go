package review_vote

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kulikovroman08/reviewlink-backend/internal/repository"
	serviceErrors "github.com/kulikovroman08/reviewlink-backend/internal/service/errors"
)

type reviewVoteService struct {
	db         *pgxpool.Pool
	voteRepo   repository.ReviewVoteRepository
	reviewRepo repository.ReviewRepository
	userRepo   repository.UserRepository
}

func NewReviewVoteService(
	db *pgxpool.Pool,
	voteRepo repository.ReviewVoteRepository,
	reviewRepo repository.ReviewRepository,
	userRepo repository.UserRepository,
) *reviewVoteService {
	return &reviewVoteService{
		db:         db,
		voteRepo:   voteRepo,
		reviewRepo: reviewRepo,
		userRepo:   userRepo,
	}
}

func (s *reviewVoteService) Vote(ctx context.Context, reviewID, userID string, value int16) error {
	if value != 1 && value != -1 {
		return serviceErrors.ErrInvalidVoteValue
	}

	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	if u.IsDeleted {
		return serviceErrors.ErrAccessDenied
	}

	_, err = s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return serviceErrors.ErrReviewNotFound
		}
		return fmt.Errorf("get review: %w", err)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Println("rollback failed:", err)
		}
	}()

	prev, err := s.voteRepo.GetVoteForUpdate(ctx, tx, reviewID, userID)
	if err != nil {
		return fmt.Errorf("get vote for update: %w", err)
	}

	var deltaHelpful int32
	var deltaUnhelpful int32

	switch {
	case prev == nil:
		if err := s.voteRepo.UpsertVote(ctx, tx, reviewID, userID, value); err != nil {
			return fmt.Errorf("upsert vote: %w", err)
		}
		if value == 1 {
			deltaHelpful = 1
		} else {
			deltaUnhelpful = 1
		}

	case *prev == value:

	default:
		if err := s.voteRepo.UpsertVote(ctx, tx, reviewID, userID, value); err != nil {
			return fmt.Errorf("upsert vote: %w", err)
		}

		if *prev == -1 && value == 1 {
			deltaUnhelpful = -1
			deltaHelpful = 1
		} else if *prev == 1 && value == -1 {
			deltaHelpful = -1
			deltaUnhelpful = 1
		}
	}

	if deltaHelpful != 0 || deltaUnhelpful != 0 {
		if err := s.voteRepo.UpdateCounters(ctx, tx, reviewID, deltaHelpful, deltaUnhelpful); err != nil {
			return fmt.Errorf("update counters: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (s *reviewVoteService) Unvote(ctx context.Context, reviewID, userID string) error {
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	if u.IsDeleted {
		return serviceErrors.ErrAccessDenied
	}

	_, err = s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return serviceErrors.ErrReviewNotFound
		}
		return fmt.Errorf("get review: %w", err)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Println("rollback failed:", err)
		}
	}()

	prev, err := s.voteRepo.GetVoteForUpdate(ctx, tx, reviewID, userID)
	if err != nil {
		return fmt.Errorf("get vote for update: %w", err)
	}

	if prev == nil {
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit tx: %w", err)
		}
		return nil
	}

	deleted, err := s.voteRepo.DeleteVote(ctx, tx, reviewID, userID)
	if err != nil {
		return fmt.Errorf("delete vote: %w", err)
	}
	if !deleted {
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit tx: %w", err)
		}
		return nil
	}

	var deltaHelpful int32
	var deltaUnhelpful int32
	if *prev == 1 {
		deltaHelpful = -1
	} else if *prev == -1 {
		deltaUnhelpful = -1
	}

	if err := s.voteRepo.UpdateCounters(ctx, tx, reviewID, deltaHelpful, deltaUnhelpful); err != nil {
		return fmt.Errorf("update counters: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
