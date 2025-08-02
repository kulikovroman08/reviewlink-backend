package controller

import (
	"github.com/kulikovroman08/reviewlink-backend/internal/service"
)

type Application struct {
	UserService  service.UserService
	PlaceService service.PlaceService
}

func NewApplication(user service.UserService, place service.PlaceService) *Application {
	return &Application{
		UserService:  user,
		PlaceService: place,
	}
}
