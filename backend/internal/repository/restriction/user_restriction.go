package restriction

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

type PostgresUserRestrictionRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRestrictionRepository(db *pgxpool.Pool) *PostgresUserRestrictionRepository {
	return &PostgresUserRestrictionRepository{db: db}
}

func (r *PostgresUserRestrictionRepository) HasActiveRestriction(
	ctx context.Context,
	userID, restrictionType string,
) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM user_restrictions 
			WHERE user_id = $1 
			AND restriction_type = $2 
			AND expires_at > $3
		)`

	var exists bool
	err := r.db.QueryRow(ctx, query, userID, restrictionType, time.Now()).Scan(&exists)
	return exists, err
}

func (r *PostgresUserRestrictionRepository) CreateRestriction(
	ctx context.Context,
	restriction *model.UserRestriction,
) error {
	query := `
		INSERT INTO user_restrictions (id, user_id, restriction_type, reason, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, restriction_type) DO NOTHING`

	_, err := r.db.Exec(
		ctx, query,
		restriction.ID, restriction.UserID, restriction.RestrictionType,
		restriction.Reason, restriction.CreatedAt, restriction.ExpiresAt,
	)
	return err
}
