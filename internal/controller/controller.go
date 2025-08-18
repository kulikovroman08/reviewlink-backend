package controller

import (
	"github.com/kulikovroman08/reviewlink-backend/internal/service"
)

type Application struct {
	Service *service.Service
}

func NewApplication(svc *service.Service) *Application {
	return &Application{
		Service: svc,
	}
}
