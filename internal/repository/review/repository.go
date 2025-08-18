package review

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

const (
	reviewTokenTable     = "review_tokens"
	reviewTokenIDColumn  = "id"
	reviewTokenPlaceID   = "place_id"
	reviewTokenValue     = "token_value"
	reviewTokenIsUsed    = "is_used"
	reviewTokenExpiresAt = "expires_at"

	reviewTable     = "reviews"
	reviewIDColumn  = "id"
	reviewUserID    = "user_id"
	reviewPlaceID   = "place_id"
	reviewTokenID   = "token_id"
	reviewContent   = "content"
	reviewRating    = "rating"
	reviewCreatedAt = "created_at"
)

type PostgresReviewRepository struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewPostgresReviewRepository(db *pgxpool.Pool) *PostgresReviewRepository {
	return &PostgresReviewRepository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *PostgresReviewRepository) GetReviewToken(ctx context.Context, token string) (*model.ReviewToken, error) {
	query, args, err := r.builder.
		Select(
			reviewTokenIDColumn,
			reviewTokenPlaceID,
			reviewTokenValue,
			reviewTokenIsUsed,
			reviewTokenExpiresAt,
		).
		From(reviewTokenTable).
		Where(sq.Eq{reviewTokenValue: token}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build GetReviewToken query: %w", err)
	}

	row := r.db.QueryRow(ctx, query, args...)

	var rt model.ReviewToken

	err = row.Scan(
		&rt.ID,
		&rt.PlaceID,
		&rt.Token,
		&rt.IsUsed,
		&rt.ExpiresAt,
	)

	if err != nil {
		return nil, fmt.Errorf("scan GetReviewToken row: %w", err)
	}

	return &rt, nil
}

func (r *PostgresReviewRepository) MarkReviewTokenUsed(ctx context.Context, tokenID uuid.UUID) error {
	query, args, err := r.builder.
		Update(reviewTokenTable).
		Set(reviewTokenIsUsed, true).
		Where(sq.Eq{reviewTokenIDColumn: tokenID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("build MarkReviewTokenUsed query: %w", err)
	}

	if _, err := r.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("exec MarkReviewTokenUsed: %w", err)
	}

	return nil
}

func (r *PostgresReviewRepository) CreateReview(ctx context.Context, review model.Review) error {
	query, args, err := r.builder.
		Insert(reviewTable).
		Columns(
			reviewIDColumn,
			reviewUserID,
			reviewPlaceID,
			reviewTokenID,
			reviewContent,
			reviewRating,
			reviewCreatedAt,
		).
		Values(
			review.ID,
			review.UserID,
			review.PlaceID,
			review.TokenID,
			review.Content,
			review.Rating,
			time.Now().UTC(),
		).
		ToSql()

	if err != nil {
		return fmt.Errorf("build CreateReview query: %w", err)
	}

	if _, err := r.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("exec CreateReview: %w", err)
	}

	return nil
}

func (r *PostgresReviewRepository) HasReviewToday(ctx context.Context, userID, placeID string) (bool, error) {
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	query, args, err := r.builder.
		Select("1").
		From(reviewTable).
		Where(sq.Eq{
			reviewUserID:  userID,
			reviewPlaceID: placeID,
		}).
		Where(sq.GtOrEq{"created_at": today}).
		Where(sq.Lt{"created_at": tomorrow}).
		Limit(1).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("build HasReviewToday query: %w", err)
	}

	row := r.db.QueryRow(ctx, query, args...)
	var dummy int
	err = row.Scan(&dummy)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("exec HasReviewToday: %w", err)
	}
	return true, nil
}
