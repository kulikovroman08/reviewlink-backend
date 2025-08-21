package place

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	placeRepo "github.com/kulikovroman08/reviewlink-backend/internal/repository/place"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

type placeService struct {
	placeRepo placeRepo.PlaceRepository
}

func NewPlaceService(placeRepo placeRepo.PlaceRepository) PlaceService {
	return &placeService{placeRepo: placeRepo}
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

	return &place, nil
}
