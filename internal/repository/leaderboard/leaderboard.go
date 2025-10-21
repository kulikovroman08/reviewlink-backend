package leaderboard

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

type Repository struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *Repository) GetTopUsers(ctx context.Context, limit int, filter model.LeaderboardFilter) ([]model.LeaderboardEntry, error) {
	builder := r.builder.
		Select(
			"users.id",
			"users.name",
			"COUNT(reviews.id) AS reviews_count",
			"ROUND(AVG(reviews.rating), 2) AS avg_rating",
		).
		From("reviews").
		Join("users ON users.id = reviews.user_id").
		Where(sq.Eq{
			"reviews.is_deleted": false,
		}).
		GroupBy("users.id", "users.name")

	having := sq.And{}
	if filter.MinRating > 0 {
		having = append(having, sq.GtOrEq{"AVG(reviews.rating)": filter.MinRating})
	}
	if filter.MinReviews > 0 {
		having = append(having, sq.GtOrEq{"COUNT(reviews.id)": filter.MinReviews})
	}
	if len(having) > 0 {
		builder = builder.Having(having)
	}

	switch filter.SortBy {
	case "rating":
		builder = builder.OrderBy("avg_rating DESC", "reviews_count DESC")
	default:
		builder = builder.OrderBy("reviews_count DESC", "avg_rating DESC")
	}

	if limit > 0 {
		builder = builder.Limit(uint64(limit))
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build SQL query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}
	defer rows.Close()

	var result []model.LeaderboardEntry
	for rows.Next() {
		var entry model.LeaderboardEntry
		if err := rows.Scan(
			&entry.ID,
			&entry.Name,
			&entry.ReviewsCount,
			&entry.AvgRating,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		result = append(result, entry)
	}

	return result, rows.Err()
}

func (r *Repository) GetTopPlaces(ctx context.Context, limit int, filter model.LeaderboardFilter) ([]model.LeaderboardEntry, error) {
	builder := r.builder.
		Select(
			"places.id",
			"places.name",
			"COUNT(reviews.id) AS reviews_count",
			"ROUND(AVG(reviews.rating), 2) AS avg_rating",
		).
		From("reviews").
		Join("places ON places.id = reviews.place_id").
		Where(sq.Eq{
			"reviews.is_deleted": false,
		}).
		GroupBy("places.id", "places.name")

	having := sq.And{}
	if filter.MinRating > 0 {
		having = append(having, sq.GtOrEq{"AVG(reviews.rating)": filter.MinRating})
	}
	if filter.MinReviews > 0 {
		having = append(having, sq.GtOrEq{"COUNT(reviews.id)": filter.MinReviews})
	}
	if len(having) > 0 {
		builder = builder.Having(having)
	}

	switch filter.SortBy {
	case "rating":
		builder = builder.OrderBy("avg_rating DESC", "reviews_count DESC")
	default:
		builder = builder.OrderBy("reviews_count DESC", "avg_rating DESC")
	}

	if limit > 0 {
		builder = builder.Limit(uint64(limit))
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build SQL query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}
	defer rows.Close()

	var result []model.LeaderboardEntry
	for rows.Next() {
		var entry model.LeaderboardEntry
		if err := rows.Scan(
			&entry.ID,
			&entry.Name,
			&entry.ReviewsCount,
			&entry.AvgRating,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		result = append(result, entry)
	}

	return result, rows.Err()
}
