package review_vote

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

const (
	reviewVoteTable = "review_votes"

	reviewVoteReviewID = "review_id"
	reviewVoteUserID   = "user_id"
	reviewVoteValue    = "value"

	reviewTable             = "reviews"
	reviewHelpfulCountCol   = "helpful_count"
	reviewUnhelpfulCountCol = "unhelpful_count"
)

type PostgresReviewVoteRepository struct {
	builder sq.StatementBuilderType
}

func NewPostgresReviewVoteRepository() *PostgresReviewVoteRepository {
	return &PostgresReviewVoteRepository{
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *PostgresReviewVoteRepository) GetVoteForUpdate(
	ctx context.Context,
	tx pgx.Tx,
	reviewID string,
	userID string,
) (*int16, error) {

	query, args, err := r.builder.
		Select(reviewVoteValue).
		From(reviewVoteTable).
		Where(sq.Eq{
			reviewVoteReviewID: reviewID,
			reviewVoteUserID:   userID,
		}).
		Suffix("FOR UPDATE").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build GetVoteForUpdate query: %w", err)
	}

	row := tx.QueryRow(ctx, query, args...)

	var value int16
	if err := row.Scan(&value); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan GetVoteForUpdate row: %w", err)
	}

	return &value, nil
}

func (r *PostgresReviewVoteRepository) UpsertVote(
	ctx context.Context,
	tx pgx.Tx,
	reviewID string,
	userID string,
	value int16,
) error {

	query := `
		INSERT INTO review_votes (review_id, user_id, value)
		VALUES ($1, $2, $3)
		ON CONFLICT (review_id, user_id)
		DO UPDATE SET
			value = EXCLUDED.value,
			updated_at = now()
	`

	if _, err := tx.Exec(ctx, query, reviewID, userID, value); err != nil {
		return fmt.Errorf("exec UpsertVote: %w", err)
	}

	return nil
}

func (r *PostgresReviewVoteRepository) DeleteVote(
	ctx context.Context,
	tx pgx.Tx,
	reviewID string,
	userID string,
) (bool, error) {

	query, args, err := r.builder.
		Delete(reviewVoteTable).
		Where(sq.Eq{
			reviewVoteReviewID: reviewID,
			reviewVoteUserID:   userID,
		}).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("build DeleteVote query: %w", err)
	}

	cmd, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("exec DeleteVote: %w", err)
	}

	return cmd.RowsAffected() > 0, nil
}

func (r *PostgresReviewVoteRepository) UpdateCounters(
	ctx context.Context,
	tx pgx.Tx,
	reviewID string,
	deltaHelpful int32,
	deltaUnhelpful int32,
) error {

	query, args, err := r.builder.
		Update(reviewTable).
		Set(reviewHelpfulCountCol,
			sq.Expr(fmt.Sprintf("%s + ?", reviewHelpfulCountCol), deltaHelpful),
		).
		Set(reviewUnhelpfulCountCol,
			sq.Expr(fmt.Sprintf("%s + ?", reviewUnhelpfulCountCol), deltaUnhelpful),
		).
		Where(sq.Eq{
			"id": reviewID,
		}).
		ToSql()

	if err != nil {
		return fmt.Errorf("build UpdateCounters query: %w", err)
	}

	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("exec UpdateCounters: %w", err)
	}

	return nil
}
