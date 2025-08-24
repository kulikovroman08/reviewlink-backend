package token

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kulikovroman08/reviewlink-backend/internal/model"
	repo "github.com/kulikovroman08/reviewlink-backend/internal/repository"
)

type Service struct {
	repo repo.TokenRepository
}

func NewService(repo repo.TokenRepository) *Service {
	return &Service{repo: repo}
}
func (s *Service) GenerateTokens(ctx context.Context, placeID string, count int) (model.GenerateTokensResult, error) {
	id, err := uuid.Parse(placeID)
	if err != nil {
		return model.GenerateTokensResult{}, fmt.Errorf("invalid place_id: %w", err)
	}

	tokens := make([]model.ReviewToken, 0, count)
	values := make([]string, 0, count)

	for i := 0; i < count; i++ {
		token := model.ReviewToken{
			ID:        uuid.New(),
			PlaceID:   id,
			Token:     uuid.New().String()[:8],
			IsUsed:    false,
			ExpiresAt: time.Now().Add(72 * time.Hour),
		}
		tokens = append(tokens, token)
		values = append(values, token.Token)
	}

	if err := s.repo.CreateTokens(ctx, tokens); err != nil {
		return model.GenerateTokensResult{}, fmt.Errorf("generate tokens: %w", err)
	}

	return model.GenerateTokensResult{Tokens: values}, nil
}
