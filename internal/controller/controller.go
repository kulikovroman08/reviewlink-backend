package controller

import (
	"github.com/kulikovroman08/reviewlink-backend/internal/service"
)

type Application struct {
	UserService   service.UserService
	PlaceService  service.PlaceService
	ReviewService service.ReviewService
}

func NewApplication(user service.UserService, place service.PlaceService, review service.ReviewService) *Application {
	return &Application{
		UserService:   user,
		PlaceService:  place,
		ReviewService: review,
	}
}
