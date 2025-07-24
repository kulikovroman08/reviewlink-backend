package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/kulikovroman08/reviewlink-backend/internal/service/user"

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
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	query, args, err := sq.
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
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build FindByEmail query: %w", err)
	}

	var u user.User
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.Points,
		&u.CreatedAt,
		&u.IsDeleted,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan FindByEmail: %w", err)
	}
	return &u, nil
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *user.User) error {
	query, args, err := sq.
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
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("build CreateUser query: %w", err)
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec CreateUser: %w", err)
	}
	return nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	query, args, err := sq.
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
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build FindByID query: %w", err)
	}

	var u user.User
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.Points,
		&u.CreatedAt,
		&u.IsDeleted,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan FindByID: %w", err)
	}

	return &u, nil
}
