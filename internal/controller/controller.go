package controller

import (
	"github.com/kulikovroman08/reviewlink-backend/internal/service"
)

type Application struct {
	UserService   service.UserService
	PlaceService  service.PlaceService
	ReviewService service.ReviewService
	TokenService  service.TokenService
	AdminService  service.AdminService
}

func NewApplication(
	user service.UserService,
	place service.PlaceService,
	review service.ReviewService,
	token service.TokenService,
	admin service.AdminService,
) *Application {
	return &Application{
		UserService:   user,
		PlaceService:  place,
		ReviewService: review,
		TokenService:  token,
		AdminService:  admin,
	}
}
