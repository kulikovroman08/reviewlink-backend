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

func (r *PostgresPlaceRepository) CreatePlace(ctx context.Context, place *model.Place) error {
	query, args, err := r.builder.
		Insert(placeTable).
		Columns(
			placeIDColumn,
			placeNameColumn,
			placeAddressColumn,
		).
		Values(
			place.ID,
			place.Name,
			place.Address,
		).
		Suffix("RETURNING created_at, is_deleted").
		ToSql()

	if err != nil {
		return fmt.Errorf("build CreatePlace query: %w", err)
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&place.CreatedAt, &place.IsDeleted)
	if err != nil {
		return fmt.Errorf("exec CreatePlace: %w", err)
	}

	return nil
}
