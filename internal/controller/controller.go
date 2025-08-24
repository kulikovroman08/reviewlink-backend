package controller

import (
	"github.com/kulikovroman08/reviewlink-backend/internal/service"
)

type Application struct {
	UserService   service.UserService
	PlaceService  service.PlaceService
	ReviewService service.ReviewService
}

func NewApplication(
	userSvc service.UserService,
	placeSvc service.PlaceService,
	reviewSvc service.ReviewService,
) *Application {
	return &Application{
		UserService:   userSvc,
		PlaceService:  placeSvc,
		ReviewService: reviewSvc,
	}
}
