package user

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kulikovroman08/reviewlink-backend/internal/service/user/model"
	"github.com/kulikovroman08/reviewlink-backend/internal/service/user/model/claims"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo UserRepository
}

func NewService(repo UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Signup(ctx context.Context, name, email, password string) (string, error) {
	existing, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("check existing user: %w", err)
	}
	if existing != nil {
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
		return "", fmt.Errorf("check existing user: %w", err)
	}
	if user == nil {
		return "", fmt.Errorf("user not found")
	}
	if user.IsDeleted {
		return "", fmt.Errorf("user is deleted")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid password")
	}
	return s.generateJWT(user)
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
