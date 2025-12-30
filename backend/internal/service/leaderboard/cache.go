package leaderboard

import (
	"context"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

type LeaderboardCache interface {
	GetUsers(ctx context.Context, limit int, filter model.LeaderboardFilter) ([]model.LeaderboardEntry, bool)
	SetUsers(ctx context.Context, limit int, filter model.LeaderboardFilter, entries []model.LeaderboardEntry)

	GetPlaces(ctx context.Context, limit int, filter model.LeaderboardFilter) ([]model.LeaderboardEntry, bool)
	SetPlaces(ctx context.Context, limit int, filter model.LeaderboardFilter, entries []model.LeaderboardEntry)

	GetBonuses(ctx context.Context) ([]model.BonusLeaderboardEntry, bool)
	SetBonuses(ctx context.Context, entries []model.BonusLeaderboardEntry)
}
