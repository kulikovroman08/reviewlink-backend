package token

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kulikovroman08/reviewlink-backend/configs"
	"github.com/kulikovroman08/reviewlink-backend/internal/model"
	repo "github.com/kulikovroman08/reviewlink-backend/internal/repository"
)

type Service struct {
	repo repo.TokenRepository
	cfg  *configs.Config
}

func NewTokenService(repo repo.TokenRepository, cfg *configs.Config) *Service {
	return &Service{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *Service) GenerateTokens(ctx context.Context, placeID string, count int) (*model.GenerateTokensResult, error) {

	tokens, values, err := generateTokens(placeID, count)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	if err := s.repo.CreateTokens(ctx, tokens); err != nil {
		return nil, fmt.Errorf("create tokens: %w", err)
	}

	return &model.GenerateTokensResult{Tokens: values}, nil
}

func (s *Service) CheckAndRefillTokens(ctx context.Context, placeID string) error {
	activeCount, err := s.repo.CountActiveTokens(ctx, placeID)
	if err != nil {
		return fmt.Errorf("count active tokens: %w", err)
	}

	if activeCount <= s.cfg.TokensThreshold {
		log.Printf("[auto-refill] active tokens for %s = %d, generating +%d",
			placeID, activeCount, s.cfg.TokensBatchSize)

		if _, err := s.GenerateTokens(ctx, placeID, s.cfg.TokensBatchSize); err != nil {
			log.Printf("[auto-refill] failed for %s: %v", placeID, err)
			return err
		}
	}

	return nil
}

func generateTokens(placeIDStr string, count int) ([]model.ReviewToken, []string, error) {
	placeID, err := uuid.Parse(placeIDStr)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid place_id: %w", err)
	}

	tokens := make([]model.ReviewToken, 0, count)
	values := make([]string, 0, count)

	for i := 0; i < count; i++ {
		token := model.ReviewToken{
			ID:        uuid.New(),
			PlaceID:   placeID,
			Token:     uuid.New().String()[:8],
			IsUsed:    false,
			ExpiresAt: time.Now().Add(72 * time.Hour),
		}
		tokens = append(tokens, token)
		values = append(values, token.Token)
	}
	return tokens, values, nil
}
