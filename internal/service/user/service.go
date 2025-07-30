package user

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
	"github.com/kulikovroman08/reviewlink-backend/internal/model"
	"github.com/kulikovroman08/reviewlink-backend/internal/model/claims"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo UserRepository
}

func NewService(repo UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetMe(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user.IsDeleted {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *Service) Signup(ctx context.Context, name, email, password string) (string, error) {
	existing, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("check existing user: %w", err)
		}
	} else if existing != nil {
		return "", fmt.Errorf("email already used")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	user := &model.User{
		ID:           uuid.NewString(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         "user",
		Points:       0,
		CreatedAt:    time.Now(),
		IsDeleted:    false,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return "", fmt.Errorf("create user: %w", err)
	}
	return s.generateJWT(user)
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("user not found: %w", err)
		}
		return "", fmt.Errorf("check existing user: %w", err)
	}
	if user.IsDeleted {
		return "", fmt.Errorf("user is deleted")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}
	return s.generateJWT(user)
}

func (s *Service) UpdateMe(ctx context.Context, req dto.UpdateUserRequest) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("find user: %w", err)
	}
	if user.IsDeleted {
		return nil, fmt.Errorf("user is deleted")
	}

	if s.shouldUpdateEmail(req.Email, user.Email) {
		var existing *model.User
		existing, err := s.repo.FindByEmail(ctx, *req.Email)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
			} else {
				return nil, fmt.Errorf("check email: %w", err)
			}
		} else if existing.ID != user.ID {
			return nil, fmt.Errorf("email already used")
		}
		user.Email = *req.Email
	}

	if req.Name != nil {
		user.Name = *req.Name
	}

	if req.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("hash password: %w", err)
		}
		user.PasswordHash = string(hash)
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}

func (s *Service) DeleteMe(ctx context.Context, userID string) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("user not found: %w", err)
		}
		return fmt.Errorf("find user: %w", err)
	}
	if user.IsDeleted {
		return fmt.Errorf("user is deleted")
	}

	if err := s.repo.SoftDeleteUser(ctx, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("user not found: %w", err)
		}
		return fmt.Errorf("soft delete user: %w", err)
	}

	return nil
}

func (s *Service) generateJWT(user *model.User) (string, error) {
	claims := claims.Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (s *Service) shouldUpdateEmail(newEmail *string, currentEmail string) bool {
	if newEmail == nil {
		return false
	}
	if *newEmail == currentEmail {
		return false
	}
	return true
}
