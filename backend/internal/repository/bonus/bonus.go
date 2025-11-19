package bonus

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
	srvErrors "github.com/kulikovroman08/reviewlink-backend/internal/service/errors"
)

const (
	bonusTable             = "bonus_rewards"
	bonusIDColumn          = "id"
	bonusUserIDColumn      = "user_id"
	bonusPlaceIDColumn     = "place_id"
	bonusRequiredPtsColumn = "required_points"
	bonusRewardTypeColumn  = "reward_type"
	bonusQRTokenColumn     = "qr_token"
	bonusIsUsedColumn      = "is_used"
	bonusUsedAtColumn      = "used_at"
)

type PostgresBonusRepository struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewPostgresBonusRepository(db *pgxpool.Pool) *PostgresBonusRepository {
	return &PostgresBonusRepository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *PostgresBonusRepository) CreateBonus(ctx context.Context, bonus *model.BonusReward) error {
	query, args, err := r.builder.
		Insert(bonusTable).
		Columns(
			bonusIDColumn,
			bonusUserIDColumn,
			bonusPlaceIDColumn,
			bonusRequiredPtsColumn,
			bonusRewardTypeColumn,
			bonusQRTokenColumn,
			bonusIsUsedColumn,
			bonusUsedAtColumn,
		).
		Values(
			bonus.ID,
			bonus.UserID,
			bonus.PlaceID,
			bonus.RequiredPoints,
			bonus.RewardType,
			bonus.QRToken,
			bonus.IsUsed,
			bonus.UsedAt,
		).
		ToSql()

	if err != nil {
		return fmt.Errorf("build CreateBonus query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec CreateBonus: %w", err)
	}

	return nil
}

func (r *PostgresBonusRepository) GetBonusesByUser(ctx context.Context, userID string) ([]model.BonusReward, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	query, args, err := r.builder.
		Select(
			bonusIDColumn,
			bonusUserIDColumn,
			bonusPlaceIDColumn,
			bonusRequiredPtsColumn,
			bonusRewardTypeColumn,
			bonusQRTokenColumn,
			bonusIsUsedColumn,
			bonusUsedAtColumn,
		).
		From(bonusTable).
		Where(sq.Eq{bonusUserIDColumn: uid}).
		OrderBy(fmt.Sprintf("%s DESC NULLS LAST", bonusUsedAtColumn)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build GetBonusesByUser query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec GetBonusesByUser: %w", err)
	}
	defer rows.Close()

	var bonuses []model.BonusReward
	for rows.Next() {
		var b model.BonusReward
		if err := rows.Scan(
			&b.ID,
			&b.UserID,
			&b.PlaceID,
			&b.RequiredPoints,
			&b.RewardType,
			&b.QRToken,
			&b.IsUsed,
			&b.UsedAt,
		); err != nil {
			return nil, fmt.Errorf("scan GetBonusesByUser: %w", err)
		}
		bonuses = append(bonuses, b)
	}

	return bonuses, rows.Err()
}

func (r *PostgresBonusRepository) MarkBonusUsed(ctx context.Context, qrToken string) error {
	query, args, err := r.builder.
		Update(bonusTable).
		Set(bonusIsUsedColumn, true).
		Set(bonusUsedAtColumn, time.Now()).
		Where(sq.Eq{bonusQRTokenColumn: qrToken}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build MarkBonusUsed query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec MarkBonusUsed: %w", err)
	}

	return nil
}

func (r *PostgresBonusRepository) GetByQRToken(ctx context.Context, qrToken string) (*model.BonusReward, error) {
	query, args, err := r.builder.
		Select(
			bonusIDColumn,
			bonusUserIDColumn,
			bonusPlaceIDColumn,
			bonusRequiredPtsColumn,
			bonusRewardTypeColumn,
			bonusQRTokenColumn,
			bonusIsUsedColumn,
			bonusUsedAtColumn,
		).
		From(bonusTable).
		Where(sq.Eq{bonusQRTokenColumn: qrToken}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build GetByQRToken query: %w", err)
	}

	row := r.db.QueryRow(ctx, query, args...)

	var b model.BonusReward
	err = row.Scan(
		&b.ID,
		&b.UserID,
		&b.PlaceID,
		&b.RequiredPoints,
		&b.RewardType,
		&b.QRToken,
		&b.IsUsed,
		&b.UsedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, srvErrors.ErrBonusNotFound
		}
		return nil, fmt.Errorf("scan GetByQRToken: %w", err)
	}

	return &b, nil
}
