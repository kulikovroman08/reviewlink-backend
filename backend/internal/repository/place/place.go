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
	placeOwnerIDColumn   = "owner_id"
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
			placeOwnerIDColumn,
			placeNameColumn,
			placeAddressColumn,
		).
		Values(
			place.ID,
			place.OwnerID,
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
			placeOwnerIDColumn,
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
	if err := row.Scan(&p.ID, &p.OwnerID, &p.Name, &p.Address, &p.CreatedAt, &p.IsDeleted); err != nil {
		return nil, err
	}

	return p, nil
}

func (r *PostgresPlaceRepository) GetPlacesByOwner(ctx context.Context, ownerID string) ([]model.Place, error) {
	ownerUUID, err := uuid.Parse(ownerID)
	if err != nil {
		return nil, fmt.Errorf("invalid owner id: %w", err)
	}

	query, args, err := r.builder.
		Select(
			placeIDColumn,
			placeOwnerIDColumn,
			placeNameColumn,
			placeAddressColumn,
			placeCreatedAtColumn,
			placeIsDeletedColumn,
		).
		From(placeTable).
		Where(sq.Eq{
			placeOwnerIDColumn:   ownerUUID,
			placeIsDeletedColumn: false,
		}).
		OrderBy(placeNameColumn + " ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build GetPlacesByOwner query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec GetPlacesByOwner: %w", err)
	}
	defer rows.Close()

	places := make([]model.Place, 0)
	for rows.Next() {
		var p model.Place
		if err := rows.Scan(
			&p.ID,
			&p.OwnerID,
			&p.Name,
			&p.Address,
			&p.CreatedAt,
			&p.IsDeleted,
		); err != nil {
			return nil, fmt.Errorf("scan GetPlacesByOwner: %w", err)
		}
		places = append(places, p)
	}

	return places, nil
}

func (r *PostgresPlaceRepository) IsOwner(ctx context.Context, placeID string, ownerID string) (bool, error) {
	placeUUID, err := uuid.Parse(placeID)
	if err != nil {
		return false, fmt.Errorf("invalid place id: %w", err)
	}
	ownerUUID, err := uuid.Parse(ownerID)
	if err != nil {
		return false, fmt.Errorf("invalid owner id: %w", err)
	}

	const q = `
		SELECT EXISTS(
			SELECT 1
			FROM places
			WHERE id = $1 AND owner_id = $2 AND is_deleted = false
		)
	`

	var ok bool
	if err := r.db.QueryRow(ctx, q, placeUUID, ownerUUID).Scan(&ok); err != nil {
		return false, fmt.Errorf("exec IsOwner: %w", err)
	}
	return ok, nil
}
