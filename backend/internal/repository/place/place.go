package place

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

const (
	placeTable           = "places"
	placeIDColumn        = "id"
	placeNameColumn      = "name"
	placeAddressColumn   = "address"
	placeCreatedAtColumn = "created_at"
	placeIsDeletedColumn = "is_deleted"
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

func (r *PostgresPlaceRepository) GetByID(ctx context.Context, placeID string) (*model.Place, error) {
	uid, err := uuid.Parse(placeID)
	if err != nil {
		return nil, fmt.Errorf("invalid place id: %w", err)
	}

	query, args, err := r.builder.
		Select(
			placeIDColumn,
			placeNameColumn,
			placeAddressColumn,
			placeCreatedAtColumn,
			placeIsDeletedColumn,
		).
		From(placeTable).
		Where(sq.Eq{
			placeIDColumn:        uid,
			placeIsDeletedColumn: false,
		}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build GetByID query: %w", err)
	}
	row := r.db.QueryRow(ctx, query, args...)

	p := new(model.Place)
	if err := row.Scan(&p.ID, &p.Name, &p.Address, &p.CreatedAt, &p.IsDeleted); err != nil {
		return nil, err
	}

	return p, nil
}
