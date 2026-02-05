package owner_request

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

const (
	ownerRequestTable = "owner_requests"

	ownerRequestIDColumn         = "id"
	ownerRequestUserIDColumn     = "user_id"
	ownerRequestSourceColumn     = "source"
	ownerRequestSourceIDColumn   = "source_id"
	ownerRequestNameColumn       = "name"
	ownerRequestStatusColumn     = "status"
	ownerRequestCreatedAtColumn  = "created_at"
	ownerRequestReviewedAtColumn = "reviewed_at"
	ownerRequestReviewedByColumn = "reviewed_by"
	ownerRequestCommentColumn    = "comment"
)

type PostgresOwnerRequestRepository struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewPostgresOwnerRequestRepository(db *pgxpool.Pool) *PostgresOwnerRequestRepository {
	return &PostgresOwnerRequestRepository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *PostgresOwnerRequestRepository) Create(ctx context.Context, req *model.OwnerRequest) error {
	query, args, err := r.builder.
		Insert(ownerRequestTable).
		Columns(
			ownerRequestIDColumn,
			ownerRequestUserIDColumn,
			ownerRequestSourceColumn,
			ownerRequestSourceIDColumn,
			ownerRequestNameColumn,
			ownerRequestStatusColumn,
		).
		Values(
			req.ID,
			req.UserID,
			req.Source,
			req.SourceID,
			req.Name,
			req.Status,
		).
		Suffix("RETURNING created_at").
		ToSql()
	if err != nil {
		return fmt.Errorf("build Create owner_request query: %w", err)
	}

	if err := r.db.QueryRow(ctx, query, args...).Scan(&req.CreatedAt); err != nil {
		return fmt.Errorf("exec Create owner_request: %w", err)
	}

	return nil
}

func (r *PostgresOwnerRequestRepository) GetByID(ctx context.Context, id string) (*model.OwnerRequest, error) {
	reqID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("parse owner_request id: %w", err)
	}

	query, args, err := r.builder.
		Select(
			ownerRequestIDColumn,
			ownerRequestUserIDColumn,
			ownerRequestSourceColumn,
			ownerRequestSourceIDColumn,
			ownerRequestNameColumn,
			ownerRequestStatusColumn,
			ownerRequestCreatedAtColumn,
			ownerRequestReviewedAtColumn,
			ownerRequestReviewedByColumn,
			ownerRequestCommentColumn,
		).
		From(ownerRequestTable).
		Where(sq.Eq{ownerRequestIDColumn: reqID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build GetByID owner_request query: %w", err)
	}

	var out model.OwnerRequest
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&out.ID,
		&out.UserID,
		&out.Source,
		&out.SourceID,
		&out.Name,
		&out.Status,
		&out.CreatedAt,
		&out.ReviewedAt,
		&out.ReviewedBy,
		&out.Comment,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("exec GetByID owner_request: %w", err)
	}

	return &out, nil
}

func (r *PostgresOwnerRequestRepository) ListPending(ctx context.Context, limit int) ([]model.OwnerRequest, error) {
	query, args, err := r.builder.
		Select(
			ownerRequestIDColumn,
			ownerRequestUserIDColumn,
			ownerRequestSourceColumn,
			ownerRequestSourceIDColumn,
			ownerRequestNameColumn,
			ownerRequestStatusColumn,
			ownerRequestCreatedAtColumn,
			ownerRequestReviewedAtColumn,
			ownerRequestReviewedByColumn,
			ownerRequestCommentColumn,
		).
		From(ownerRequestTable).
		Where(sq.Eq{ownerRequestStatusColumn: model.OwnerRequestPending}).
		OrderBy(ownerRequestCreatedAtColumn + " DESC").
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build ListPending owner_request query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec ListPending owner_request: %w", err)
	}
	defer rows.Close()

	var out []model.OwnerRequest
	for rows.Next() {
		var rreq model.OwnerRequest
		if err := rows.Scan(
			&rreq.ID,
			&rreq.UserID,
			&rreq.Source,
			&rreq.SourceID,
			&rreq.Name,
			&rreq.Status,
			&rreq.CreatedAt,
			&rreq.ReviewedAt,
			&rreq.ReviewedBy,
			&rreq.Comment,
		); err != nil {
			return nil, fmt.Errorf("scan ListPending owner_request: %w", err)
		}
		out = append(out, rreq)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows ListPending owner_request: %w", err)
	}

	return out, nil
}

func (r *PostgresOwnerRequestRepository) Approve(ctx context.Context, id string, reviewedBy string, comment *string) error {
	reqID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("parse owner_request id: %w", err)
	}
	adminID, err := uuid.Parse(reviewedBy)
	if err != nil {
		return fmt.Errorf("parse reviewedBy: %w", err)
	}

	query, args, err := r.builder.
		Update(ownerRequestTable).
		Set(ownerRequestStatusColumn, model.OwnerRequestApproved).
		Set(ownerRequestReviewedAtColumn, sq.Expr("now()")).
		Set(ownerRequestReviewedByColumn, adminID).
		Set(ownerRequestCommentColumn, comment).
		Where(sq.Eq{ownerRequestIDColumn: reqID}).
		Where(sq.Eq{ownerRequestStatusColumn: model.OwnerRequestPending}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build Approve owner_request query: %w", err)
	}

	ct, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec Approve owner_request: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("approve owner_request: not found or not pending")
	}

	return nil
}

func (r *PostgresOwnerRequestRepository) Reject(ctx context.Context, id string, reviewedBy string, comment *string) error {
	reqID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("parse owner_request id: %w", err)
	}
	adminID, err := uuid.Parse(reviewedBy)
	if err != nil {
		return fmt.Errorf("parse reviewedBy: %w", err)
	}

	query, args, err := r.builder.
		Update(ownerRequestTable).
		Set(ownerRequestStatusColumn, model.OwnerRequestRejected).
		Set(ownerRequestReviewedAtColumn, sq.Expr("now()")).
		Set(ownerRequestReviewedByColumn, adminID).
		Set(ownerRequestCommentColumn, comment).
		Where(sq.Eq{ownerRequestIDColumn: reqID}).
		Where(sq.Eq{ownerRequestStatusColumn: model.OwnerRequestPending}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build Reject owner_request query: %w", err)
	}

	ct, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec Reject owner_request: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("reject owner_request: not found or not pending")
	}

	return nil
}
