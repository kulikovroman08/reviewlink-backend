package review_reply

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

const (
	reviewReplyTable = "review_replies"

	colID        = "id"
	colReviewID  = "review_id"
	colAdminID   = "admin_id"
	colContent   = "content"
	colCreatedAt = "created_at"
	colUpdatedAt = "updated_at"
)

type PostgresReviewReplyRepository struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewPostgresReviewReplyRepository(db *pgxpool.Pool) *PostgresReviewReplyRepository {
	return &PostgresReviewReplyRepository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *PostgresReviewReplyRepository) CreateReply(ctx context.Context, reply *model.ReviewReply) error {
	if reply == nil {
		return fmt.Errorf("reply is nil")
	}
	if reply.ID == uuid.Nil {
		reply.ID = uuid.New()
	}
	if reply.CreatedAt.IsZero() {
		reply.CreatedAt = time.Now().UTC()
	}

	query, args, err := r.builder.
		Insert(reviewReplyTable).
		Columns(
			colID,
			colReviewID,
			colAdminID,
			colContent,
			colCreatedAt,
		).
		Values(
			reply.ID,
			reply.ReviewID,
			reply.AdminID,
			reply.Content,
			reply.CreatedAt,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("build CreateReply query: %w", err)
	}

	if _, err := r.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("exec CreateReply: %w", err)
	}

	return nil
}

func (r *PostgresReviewReplyRepository) UpdateReply(ctx context.Context, reviewID string, content string) (*model.ReviewReply, error) {
	rid, err := uuid.Parse(reviewID)
	if err != nil {
		return nil, fmt.Errorf("invalid review id: %w", err)
	}

	now := time.Now().UTC()

	query, args, err := r.builder.
		Update(reviewReplyTable).
		Set(colContent, content).
		Set(colUpdatedAt, now).
		Where(sq.Eq{colReviewID: rid}).
		Suffix(fmt.Sprintf("RETURNING %s, %s, %s, %s, %s, %s",
			colID, colReviewID, colAdminID, colContent, colCreatedAt, colUpdatedAt,
		)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build UpdateReply query: %w", err)
	}

	row := r.db.QueryRow(ctx, query, args...)

	var rr model.ReviewReply
	var updatedAt sql.NullTime

	if err := row.Scan(
		&rr.ID,
		&rr.ReviewID,
		&rr.AdminID,
		&rr.Content,
		&rr.CreatedAt,
		&updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("scan UpdateReply row: %w", err)
	}

	if updatedAt.Valid {
		rr.UpdatedAt = &updatedAt.Time
	}

	return &rr, nil
}

func (r *PostgresReviewReplyRepository) GetByReviewID(ctx context.Context, reviewID string) (*model.ReviewReply, error) {
	rid, err := uuid.Parse(reviewID)
	if err != nil {
		return nil, fmt.Errorf("invalid review id: %w", err)
	}

	query, args, err := r.builder.
		Select(
			colID,
			colReviewID,
			colAdminID,
			colContent,
			colCreatedAt,
			colUpdatedAt,
		).
		From(reviewReplyTable).
		Where(sq.Eq{colReviewID: rid}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build GetByReviewID query: %w", err)
	}

	row := r.db.QueryRow(ctx, query, args...)

	var rr model.ReviewReply
	var updatedAt sql.NullTime

	if err := row.Scan(
		&rr.ID,
		&rr.ReviewID,
		&rr.AdminID,
		&rr.Content,
		&rr.CreatedAt,
		&updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("scan GetByReviewID row: %w", err)
	}

	if updatedAt.Valid {
		rr.UpdatedAt = &updatedAt.Time
	}

	return &rr, nil
}

func (r *PostgresReviewReplyRepository) GetByReviewIDs(ctx context.Context, reviewIDs []string) (map[string]*model.ReviewReply, error) {
	if len(reviewIDs) == 0 {
		return map[string]*model.ReviewReply{}, nil
	}

	ids := make([]uuid.UUID, 0, len(reviewIDs))
	for _, s := range reviewIDs {
		uid, err := uuid.Parse(s)
		if err != nil {
			return nil, fmt.Errorf("invalid review id: %w", err)
		}
		ids = append(ids, uid)
	}

	query, args, err := r.builder.
		Select(
			colID,
			colReviewID,
			colAdminID,
			colContent,
			colCreatedAt,
			colUpdatedAt,
		).
		From(reviewReplyTable).
		Where(sq.Eq{colReviewID: ids}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build GetByReviewIDs query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec GetByReviewIDs: %w", err)
	}
	defer rows.Close()

	out := make(map[string]*model.ReviewReply, len(reviewIDs))

	for rows.Next() {
		var rr model.ReviewReply
		var updatedAt sql.NullTime

		if err := rows.Scan(
			&rr.ID,
			&rr.ReviewID,
			&rr.AdminID,
			&rr.Content,
			&rr.CreatedAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan GetByReviewIDs row: %w", err)
		}

		if updatedAt.Valid {
			rr.UpdatedAt = &updatedAt.Time
		}

		out[rr.ReviewID.String()] = &rr
	}

	return out, nil
}
