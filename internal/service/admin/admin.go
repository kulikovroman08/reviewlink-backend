package admin

import (
	"context"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
	"github.com/kulikovroman08/reviewlink-backend/internal/repository"
)

type AdminService struct {
	repo repository.AdminRepository
}

func NewAdminService(repo repository.AdminRepository) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) GetStats(ctx context.Context) (*model.AdminStats, error) {
	return s.repo.GetAdminStats(ctx)
}
