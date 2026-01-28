package place

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

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

	var owner pgtype.UUID
	var addr sql.NullString

	if err := row.Scan(&p.ID, &owner, &p.Name, &addr, &p.CreatedAt, &p.IsDeleted); err != nil {
		return nil, err
	}

	if owner.Valid {
		u, err := uuid.FromBytes(owner.Bytes[:])
		if err != nil {
			return nil, fmt.Errorf("decode owner_id uuid: %w", err)
		}
		p.OwnerID = &u
	} else {
		p.OwnerID = nil
	}

	if addr.Valid {
		p.Address = addr.String
	} else {
		p.Address = ""
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
		var owner pgtype.UUID

		if err := rows.Scan(
			&p.ID,
			&owner,
			&p.Name,
			&p.Address,
			&p.CreatedAt,
			&p.IsDeleted,
		); err != nil {
			return nil, fmt.Errorf("scan GetPlacesByOwner: %w", err)
		}

		if owner.Valid {
			u, err := uuid.FromBytes(owner.Bytes[:])
			if err != nil {
				return nil, fmt.Errorf("decode owner_id: %w", err)
			}
			p.OwnerID = &u
		} else {
			p.OwnerID = nil
		}

		places = append(places, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error GetPlacesByOwner: %w", err)
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

func (r *PostgresPlaceRepository) GetPublicMetaBySourceIDs(
	ctx context.Context,
	source string,
	sourceIDs []string,
) (map[string]model.PublicPlaceMeta, error) {

	if len(sourceIDs) == 0 {
		return map[string]model.PublicPlaceMeta{}, nil
	}

	const q = `
		SELECT
			p.source_id,
			p.id AS place_id,
			COUNT(r.id) AS reviews_count,
			AVG(r.rating)::float8 AS avg_rating,
			(p.owner_id IS NOT NULL) AS has_owner
		FROM places p
		LEFT JOIN reviews r
		  ON r.place_id = p.id
		 AND r.is_deleted = false
		WHERE p.source = $1
		  AND p.source_id = ANY($2)
		  AND p.is_deleted = false
		GROUP BY p.source_id, p.id, p.owner_id
	`

	rows, err := r.db.Query(ctx, q, source, sourceIDs)
	if err != nil {
		return nil, fmt.Errorf("exec GetPublicMetaBySourceIDs: %w", err)
	}
	defer rows.Close()

	result := make(map[string]model.PublicPlaceMeta, len(sourceIDs))

	for rows.Next() {
		var (
			sourceID     string
			placeID      uuid.UUID
			reviewsCount int
			avgRating    pgtype.Float8
			hasOwner     bool
		)

		if err := rows.Scan(
			&sourceID,
			&placeID,
			&reviewsCount,
			&avgRating,
			&hasOwner,
		); err != nil {
			return nil, fmt.Errorf("scan GetPublicMetaBySourceIDs: %w", err)
		}

		meta := model.PublicPlaceMeta{
			SourceID:     sourceID,
			PlaceID:      placeID,
			ReviewsCount: reviewsCount,
			HasOwner:     hasOwner,
		}

		if avgRating.Valid {
			meta.Rating = &avgRating.Float64
		} else {
			meta.Rating = nil
		}

		result[sourceID] = meta
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows GetPublicMetaBySourceIDs: %w", err)
	}

	return result, nil
}

func (r *PostgresPlaceRepository) EnsureFromPublic(
	ctx context.Context,
	source string,
	sourceID string,
	name string,
) (string, error) {

	const q = `
		INSERT INTO places (id, name, source, source_id, owner_id, is_deleted, created_at)
		VALUES (gen_random_uuid(), $1, $2, $3, NULL, false, NOW())
		ON CONFLICT (source, source_id)
		DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`

	var placeID uuid.UUID

	if err := r.db.QueryRow(ctx, q, name, source, sourceID).Scan(&placeID); err != nil {
		return "", fmt.Errorf("exec EnsureFromPublic: %w", err)
	}

	return placeID.String(), nil
}
