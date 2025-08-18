package service

import (
	"github.com/kulikovroman08/reviewlink-backend/internal/service/place"
	"github.com/kulikovroman08/reviewlink-backend/internal/service/review"
	"github.com/kulikovroman08/reviewlink-backend/internal/service/user"
)

type Service struct {
	UserRepo   user.UserRepository
	ReviewRepo review.ReviewRepository
	PlaceRepo  place.PlaceRepository
}

func NewService(
	userRepo user.UserRepository,
	reviewRepo review.ReviewRepository,
	placeRepo place.PlaceRepository,
) *Service {
	return &Service{
		UserRepo:   userRepo,
		ReviewRepo: reviewRepo,
		PlaceRepo:  placeRepo,
	}
}
