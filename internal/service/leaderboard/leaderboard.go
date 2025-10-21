package leaderboard

import (
	"context"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
	"github.com/kulikovroman08/reviewlink-backend/internal/repository"
)

type LeaderboardService struct {
	repo repository.LeaderboardRepository
}

func NewService(repo repository.LeaderboardRepository) *LeaderboardService {
	return &LeaderboardService{repo: repo}
}

func (s *LeaderboardService) GetUserLeaderboard(ctx context.Context, limit int, filter model.LeaderboardFilter) ([]model.LeaderboardEntry, error) {
	return s.repo.GetTopUsers(ctx, limit, filter)
}

func (s *LeaderboardService) GetPlaceLeaderboard(ctx context.Context, limit int, filter model.LeaderboardFilter) ([]model.LeaderboardEntry, error) {
	return s.repo.GetTopPlaces(ctx, limit, filter)
}
