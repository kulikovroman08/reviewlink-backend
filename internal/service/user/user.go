package user

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/kulikovroman08/reviewlink-backend/internal/repository/user"

	"github.com/jackc/pgx/v5"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kulikovroman08/reviewlink-backend/internal/model"
	"github.com/kulikovroman08/reviewlink-backend/internal/model/claims"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo user.UserRepository
}

func NewUserService(userRepo user.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUser(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}

func (s *userService) Signup(ctx context.Context, name, email, password string) (string, error) {
	existing, err := s.userRepo.FindAnyByEmail(ctx, email)
	if err != nil {
		if isUnexpectedErr(err) {
			return "", fmt.Errorf("check existing user: %w", err)
		}
		existing = nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	if existing == nil {
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

		if err := s.userRepo.CreateUser(ctx, user); err != nil {
			return "", fmt.Errorf("create user: %w", err)
		}

		return s.generateJWT(user)
	}

	if existing.IsDeleted {
		existing.Name = name
		existing.PasswordHash = string(hashedPassword)
		existing.IsDeleted = false

		if err := s.userRepo.UpdateUser(ctx, existing); err != nil {
			return "", fmt.Errorf("restore user: %w", err)
		}
		return s.generateJWT(existing)
	}

	return "", errors.New("email already used")
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("user not found: %w", err)
		}
		return "", fmt.Errorf("check existing user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	return s.generateJWT(user)
}

func (s *userService) UpdateUser(ctx context.Context, user model.User, password string) (model.User, error) {
	current, err := s.userRepo.FindByID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, fmt.Errorf("user not found: %w", err)
		}

		return model.User{}, fmt.Errorf("find user: %w", err)
	}

	if s.shouldUpdateEmail(&user.Email, current.Email) {
		existing, err := s.userRepo.FindByEmail(ctx, user.Email)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {

			} else {
				return model.User{}, fmt.Errorf("check email: %w", err)
			}
		} else if existing.ID != user.ID {
			return model.User{}, fmt.Errorf("email already used")
		}

		current.Email = user.Email
	}

	if user.Name != "" {
		current.Name = user.Name
	}

	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return model.User{}, fmt.Errorf("hash password: %w", err)
		}
		current.PasswordHash = string(hash)
	}

	if err := s.userRepo.UpdateUser(ctx, current); err != nil {
		return model.User{}, fmt.Errorf("update user: %w", err)
	}

	return *current, nil
}

func (s *userService) DeleteUser(ctx context.Context, userID string) error {
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("user not found: %w", err)
		}
		return fmt.Errorf("find user: %w", err)
	}

	if err := s.userRepo.SoftDeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("soft delete user: %w", err)
	}

	return nil
}

func (s *userService) generateJWT(user *model.User) (string, error) {
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

func (s *userService) shouldUpdateEmail(newEmail *string, currentEmail string) bool {
	if newEmail == nil {
		return false
	}
	if *newEmail == currentEmail {
		return false
	}
	return true
}

func isUnexpectedErr(err error) bool {
	return !errors.Is(err, pgx.ErrNoRows)
}
