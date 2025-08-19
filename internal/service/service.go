package service

import (
	"github.com/kulikovroman08/reviewlink-backend/internal/repository"
)

type Service struct {
	UserRepo   repository.UserRepository
	PlaceRepo  repository.PlaceRepository
	ReviewRepo repository.ReviewRepository
}

func NewService(
	userRepo repository.UserRepository,
	reviewRepo repository.ReviewRepository,
	placeRepo repository.PlaceRepository,
) *Service {
	return &Service{
		UserRepo:   userRepo,
		ReviewRepo: reviewRepo,
		PlaceRepo:  placeRepo,
	}
}
