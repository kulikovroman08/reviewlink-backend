package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/google/uuid"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	sq "github.com/Masterminds/squirrel"
)

const (
	userTable              = "users"
	userIDColumn           = "id"
	userNameColumn         = "name"
	userEmailColumn        = "email"
	userPasswordHashColumn = "password_hash"
	userRoleColumn         = "role"
	userPointsColumn       = "points"
	userCreatedAtColumn    = "created_at"
	userIsDeletedColumn    = "is_deleted"
)

type PostgresUserRepository struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	query, args, err := r.builder.
		Select(
			userIDColumn,
			userNameColumn,
			userEmailColumn,
			userPasswordHashColumn,
			userRoleColumn,
			userPointsColumn,
			userCreatedAtColumn,
			userIsDeletedColumn,
		).
		From(userTable).
		Where(sq.And{
			sq.Eq{userIDColumn: id},
			sq.Eq{userIsDeletedColumn: false},
		}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build FindByID query: %w", err)
	}

	var u model.User

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.Points,
		&u.CreatedAt,
		&u.IsDeleted,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, pgx.ErrNoRows
	}
	if err != nil {
		return nil, fmt.Errorf("scan FindByID: %w", err)
	}

	return &u, nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query, args, err := r.builder.
		Select(
			userIDColumn,
			userNameColumn,
			userEmailColumn,
			userPasswordHashColumn,
			userRoleColumn,
			userPointsColumn,
			userCreatedAtColumn,
			userIsDeletedColumn,
		).
		From(userTable).
		Where(sq.And{
			sq.Eq{userEmailColumn: email},
			sq.Eq{userIsDeletedColumn: false},
		}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build FindByEmail query: %w", err)
	}

	var u model.User

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.Points,
		&u.CreatedAt,
		&u.IsDeleted,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, pgx.ErrNoRows
	}
	if err != nil {
		return nil, fmt.Errorf("scan FindByEmail: %w", err)
	}

	return &u, nil
}

func (r *PostgresUserRepository) FindAnyByEmail(ctx context.Context, email string) (*model.User, error) {
	query, args, err := r.builder.
		Select(
			userIDColumn,
			userNameColumn,
			userEmailColumn,
			userPasswordHashColumn,
			userRoleColumn,
			userPointsColumn,
			userCreatedAtColumn,
			userIsDeletedColumn,
		).
		From(userTable).
		Where(sq.Eq{userEmailColumn: email}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build FindAnyByEmail query: %w", err)
	}

	var u model.User

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.Points,
		&u.CreatedAt,
		&u.IsDeleted,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, pgx.ErrNoRows
	}
	if err != nil {
		return nil, fmt.Errorf("scan FindAnyByEmail: %w", err)
	}

	return &u, nil
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *model.User) error {
	query, args, err := r.builder.
		Insert(userTable).
		Columns(
			userIDColumn,
			userNameColumn,
			userEmailColumn,
			userPasswordHashColumn,
			userRoleColumn,
			userPointsColumn,
			userCreatedAtColumn,
			userIsDeletedColumn,
		).
		Values(
			user.ID,
			user.Name,
			user.Email,
			user.PasswordHash,
			user.Role,
			user.Points,
			user.CreatedAt,
			user.IsDeleted,
		).
		ToSql()

	if err != nil {
		return fmt.Errorf("build CreateUser query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec CreateUser: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	query, args, err := r.builder.
		Update(userTable).
		Set(userNameColumn, user.Name).
		Set(userEmailColumn, user.Email).
		Set(userPasswordHashColumn, user.PasswordHash).
		Set(userIsDeletedColumn, user.IsDeleted).
		Where(sq.Eq{userIDColumn: user.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("build UpdateUser query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec UpdateUser: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) SoftDeleteUser(ctx context.Context, id string) error {
	query, args, err := r.builder.
		Update(userTable).
		Set(userIsDeletedColumn, true).
		Where(sq.Eq{userIDColumn: id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("build SoftDeleteUser query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)

	if err != nil {
		return fmt.Errorf("exec SoftDeleteUser: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) AddPoints(ctx context.Context, userID uuid.UUID, points int) error {
	query, args, err := r.builder.
		Update(userTable).
		Set(userPointsColumn, sq.Expr(fmt.Sprintf("%s + ?", userPointsColumn), points)).
		Where(sq.Eq{userIDColumn: userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build AddPoints query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23514":
				return fmt.Errorf("points limit exceeded: %w", err)
			}
		}

		return fmt.Errorf("exec AddPoints: %w", err)
	}

	return nil
}
