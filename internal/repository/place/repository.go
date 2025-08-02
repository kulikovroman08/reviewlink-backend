package place

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

const (
	placeTable         = "places"
	placeIDColumn      = "id"
	placeNameColumn    = "name"
	placeAddressColumn = "address"
	placeCreatedAt     = "created_at"
	placeIsDeleted     = "is_deleted"
)

type PostgresPlaceRepository struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewPostgresPlaceRepository(db *pgxpool.Pool) *PostgresPlaceRepository {
	return &PostgresPlaceRepository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *PostgresPlaceRepository) CreatePlace(ctx context.Context, p model.Place) error {
	query, args, err := r.builder.
		Insert(placeTable).
		Columns(
			placeIDColumn,
			placeNameColumn,
			placeAddressColumn,
			placeCreatedAt,
			placeIsDeleted,
		).
		Values(
			p.ID,
			p.Name,
			p.Address,
			p.CreatedAt,
			p.IsDeleted,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("build CreatePlace query: %w", err)
	}
	if _, err := r.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("exec CreatePlace: %w", err)
	}
	return nil
}
