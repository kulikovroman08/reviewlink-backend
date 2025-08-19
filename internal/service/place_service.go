package service

import (
	"context"
	"fmt"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

func (s *Service) CreatePlace(ctx context.Context, place model.Place) (*model.Place, error) {
	if place.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if place.Address == "" {
		return nil, fmt.Errorf("address is required")
	}

	if err := s.PlaceRepo.CreatePlace(ctx, &place); err != nil {
		return nil, fmt.Errorf("failed to create place: %w", err)
	}

	return &place, nil
}
