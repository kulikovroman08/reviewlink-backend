package controller

import (
	"github.com/kulikovroman08/reviewlink-backend/internal/service"
)

type Application struct {
	UserService service.UserService
}

func NewApplication(userService service.UserService) *Application {
	return &Application{
		UserService: userService,
	}
}
