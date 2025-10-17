package token

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

const (
	reviewTokensTable          = "review_tokens"
	reviewTokenIDColumn        = "id"
	reviewTokenPlaceIDColumn   = "place_id"
	reviewTokenValueColumn     = "token_value"
	reviewTokenIsUsedColumn    = "is_used"
	reviewTokenExpiresAtColumn = "expires_at"
)

type PostgresTokenRepository struct {
	db   *pgxpool.Pool
	psql sq.StatementBuilderType
}

func NewPostgresTokenRepository(db *pgxpool.Pool) *PostgresTokenRepository {
	return &PostgresTokenRepository{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *PostgresTokenRepository) CreateTokens(ctx context.Context, tokens []model.ReviewToken) error {
	if len(tokens) == 0 {
		return nil
	}

	builder := r.psql.Insert(reviewTokensTable).
		Columns(
			reviewTokenIDColumn,
			reviewTokenPlaceIDColumn,
			reviewTokenValueColumn,
			reviewTokenExpiresAtColumn,
			reviewTokenIsUsedColumn,
		)

	for _, t := range tokens {
		builder = builder.Values(
			t.ID,
			t.PlaceID,
			t.Token,
			t.ExpiresAt,
			t.IsUsed,
		)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("build insert tokens: %w", err)
	}

	if _, err := r.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("insert tokens: %w", err)
	}

	return nil
}

func (r *PostgresTokenRepository) CountActiveTokens(ctx context.Context, placeID string) (int, error) {
	uid, err := uuid.Parse(placeID)
	if err != nil {
		return 0, fmt.Errorf("invalid placeID: %w", err)
	}

	query, args, err := r.psql.
		Select("COUNT(*)").
		From(reviewTokensTable).
		Where(sq.Eq{
			reviewTokenPlaceIDColumn: uid,
			reviewTokenIsUsedColumn:  false,
		}).ToSql()
	if err != nil {
		return 0, fmt.Errorf("build count query: %w", err)
	}

	var count int
	if err := r.db.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("execute count query: %w", err)
	}
	return count, nil
}
