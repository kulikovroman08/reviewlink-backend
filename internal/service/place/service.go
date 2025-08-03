package place

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

type Service struct {
	repo PlaceRepository
}

func NewService(repo PlaceRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreatePlace(ctx context.Context, req dto.CreatePlaceRequest) (dto.CreatePlaceResponse, error) {
	if req.Name == "" {
		return dto.CreatePlaceResponse{}, fmt.Errorf("name is required")
	}

	if req.Address == "" {
		return dto.CreatePlaceResponse{}, fmt.Errorf("address is required")
	}

	id := uuid.New()

	newPlace := model.Place{
		ID:        id,
		Name:      req.Name,
		Address:   req.Address,
		CreatedAt: time.Now().UTC(),
		IsDeleted: false,
	}

	if err := s.repo.CreatePlace(ctx, newPlace); err != nil {
		return dto.CreatePlaceResponse{}, fmt.Errorf("failed to create place: %w", err)
	}

	return dto.CreatePlaceResponse{ID: id.String()}, nil
}
