package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

func (s *Service) CreatePlace(ctx context.Context, place model.Place) (dto.CreatePlaceResponse, error) {
	if place.Name == "" {
		return dto.CreatePlaceResponse{}, fmt.Errorf("name is required")
	}

	if place.Address == "" {
		return dto.CreatePlaceResponse{}, fmt.Errorf("address is required")
	}

	place.ID = uuid.New()
	place.CreatedAt = time.Now().UTC()
	place.IsDeleted = false

	if err := s.PlaceRepo.CreatePlace(ctx, place); err != nil {
		return dto.CreatePlaceResponse{}, fmt.Errorf("failed to create place: %w", err)
	}

	return dto.CreatePlaceResponse{ID: place.ID.String()}, nil
}
