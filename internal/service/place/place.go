package place

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kulikovroman08/reviewlink-backend/configs"
	"github.com/kulikovroman08/reviewlink-backend/internal/repository"
	"github.com/kulikovroman08/reviewlink-backend/internal/service/token"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

type placeService struct {
	placeRepo    repository.PlaceRepository
	tokenService *token.Service
	cfg          *configs.Config
}

func NewPlaceService(placeRepo repository.PlaceRepository, tokenService *token.Service, cfg *configs.Config) *placeService {
	return &placeService{
		placeRepo:    placeRepo,
		tokenService: tokenService,
		cfg:          cfg,
	}
}

func (s *placeService) CreatePlace(ctx context.Context, place model.Place) (*model.Place, error) {
	if place.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if place.Address == "" {
		return nil, fmt.Errorf("address is required")
	}

	place.ID = uuid.New()

	if err := s.placeRepo.CreatePlace(ctx, &place); err != nil {
		return nil, fmt.Errorf("failed to create place: %w", err)
	}

	count := s.cfg.TokensAutoCount
	if _, err := s.tokenService.GenerateTokens(ctx, place.ID.String(), count); err != nil {
		fmt.Printf("failed to auto-generate tokens for place %s: %v\n", place.ID, err)
	}

	return &place, nil
}
