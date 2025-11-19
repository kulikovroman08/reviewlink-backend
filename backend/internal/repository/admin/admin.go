package admin

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

type PostgresAdminRepository struct {
	db   *pgxpool.Pool
	psql sq.StatementBuilderType
}

func NewPostgresAdminRepository(db *pgxpool.Pool) *PostgresAdminRepository {
	return &PostgresAdminRepository{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *PostgresAdminRepository) GetAdminStats(ctx context.Context) (*model.AdminStats, error) {
	var stats model.AdminStats

	query := `
		SELECT
			(SELECT COUNT(*) FROM users WHERE is_deleted = false) AS total_users,
			(SELECT COUNT(*) FROM reviews WHERE is_deleted = false) AS total_reviews,
			COALESCE((SELECT AVG(rating)::float FROM reviews WHERE is_deleted = false), 0) AS average_rating,
			COALESCE((SELECT COUNT(*) FROM bonus_rewards), 0) AS total_bonuses;
	`

	if err := r.db.QueryRow(ctx, query).Scan(
		&stats.TotalUsers,
		&stats.TotalReviews,
		&stats.AverageRating,
		&stats.TotalBonuses,
	); err != nil {
		return nil, fmt.Errorf("get admin stats: %w", err)
	}

	return &stats, nil
}
