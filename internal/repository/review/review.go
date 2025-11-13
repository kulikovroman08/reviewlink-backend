package review

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"

	sq "github.com/Masterminds/squirrel"
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

	reviewTable        = "reviews"
	reviewIDColumn     = "id"
	reviewUserID       = "user_id"
	reviewPlaceID      = "place_id"
	reviewTokenID      = "token_id"
	reviewContent      = "content"
	reviewRating       = "rating"
	reviewCreatedAt    = "created_at"
	reviewIsDeletedCol = "is_deleted"
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

func (r *PostgresReviewRepository) MarkReviewTokenUsed(ctx context.Context, tokenID string) error {
	uuidID, err := uuid.Parse(tokenID)
	if err != nil {
		return fmt.Errorf("invalid token ID: %w", err)
	}

	query, args, err := r.builder.
		Update(reviewTokenTable).
		Set(reviewTokenIsUsed, true).
		Where(sq.Eq{reviewTokenIDColumn: uuidID}).
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

func (r *PostgresReviewRepository) FindReviews(ctx context.Context, placeID string, filter model.ReviewFilter) ([]model.Review, error) {
	uid, err := uuid.Parse(placeID)
	if err != nil {
		return nil, fmt.Errorf("invalid place id: %w", err)
	}

	builder := r.builder.
		Select(
			reviewIDColumn,
			reviewUserID,
			reviewPlaceID,
			reviewTokenID,
			reviewContent,
			reviewRating,
			reviewCreatedAt,
		).
		From(reviewTable).
		Where(sq.Eq{
			reviewPlaceID:      uid,
			reviewIsDeletedCol: false,
		})

	if filter.HasRating {
		builder = builder.Where(sq.Eq{reviewRating: filter.Rating})
	}
	if filter.FromDate != nil {
		builder = builder.Where(sq.GtOrEq{reviewCreatedAt: *filter.FromDate})
	}
	if filter.ToDate != nil {
		builder = builder.Where(sq.LtOrEq{reviewCreatedAt: *filter.ToDate})
	}

	switch filter.Sort {
	case "date_asc":
		builder = builder.OrderBy(reviewCreatedAt + " ASC")
	default:
		builder = builder.OrderBy(reviewCreatedAt + " DESC")
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build FindReviews query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec FindReviews: %w", err)
	}
	defer rows.Close()

	reviews := make([]model.Review, 0)
	for rows.Next() {
		var rev model.Review
		err := rows.Scan(
			&rev.ID,
			&rev.UserID,
			&rev.PlaceID,
			&rev.TokenID,
			&rev.Content,
			&rev.Rating,
			&rev.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan FindReviews row: %w", err)
		}
		reviews = append(reviews, rev)
	}
	return reviews, nil
}

func (r *PostgresReviewRepository) UpdateReview(ctx context.Context, reviewID, userID string, content string, rating int) error {
	now := time.Now()

	query, args, err := r.builder.
		Update(reviewTable).
		Set(reviewContent, content).
		Set(reviewRating, rating).
		Set("updated_at", now).
		Where(sq.Eq{
			reviewIDColumn: reviewID,
			reviewUserID:   userID,
		}).
		ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PostgresReviewRepository) DeleteReview(ctx context.Context, reviewID, userID string) error {
	query, args, err := r.builder.
		Update(reviewTable).
		Set(reviewIsDeletedCol, true).
		Where(sq.Eq{
			reviewIDColumn: reviewID,
			reviewUserID:   userID,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build DeleteReview query: %w", err)
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec DeleteReview: %w", err)
	}

	if res.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PostgresReviewRepository) CountLowRatingReviews(ctx context.Context, userID string, days int) (int, error) {
	query := `
        SELECT COUNT(*)
        FROM reviews
        WHERE user_id = $1
          AND rating = 1
          AND created_at > NOW() - INTERVAL '1 day' * $2
          AND is_deleted = false`

	var count int
	err := r.db.QueryRow(ctx, query, userID, days).Scan(&count)
	return count, err
}
