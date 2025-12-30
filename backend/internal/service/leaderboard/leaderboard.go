package leaderboard

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
	"github.com/kulikovroman08/reviewlink-backend/internal/repository"
	"golang.org/x/sync/singleflight"
)

type LeaderboardService struct {
	repo  repository.LeaderboardRepository
	cache LeaderboardCache
	sf    singleflight.Group
}

func NewService(repo repository.LeaderboardRepository, cache LeaderboardCache) *LeaderboardService {
	return &LeaderboardService{
		repo:  repo,
		cache: cache,
	}
}

func (s *LeaderboardService) GetUserLeaderboard(ctx context.Context, limit int, filter model.LeaderboardFilter) ([]model.LeaderboardEntry, error) {
	if s.cache != nil {
		if entries, ok := s.cache.GetUsers(ctx, limit, filter); ok {
			return entries, nil
		}
	}

	sfKey := fmt.Sprintf(
		"leaderboard:users:limit=%d:sort=%s:minr=%s:minrev=%d",
		limit,
		filter.SortBy,
		strconv.FormatFloat(filter.MinRating, 'f', -1, 64),
		filter.MinReviews,
	)
	v, err, _ := s.sf.Do(sfKey, func() (any, error) {
		if s.cache != nil {
			if entries, ok := s.cache.GetUsers(ctx, limit, filter); ok {
				return entries, nil
			}
		}

		log.Println("DB QUERY leaderboard/users")
		entries, err := s.repo.GetTopUsers(ctx, limit, filter)
		if err != nil {
			return nil, err
		}

		if s.cache != nil {
			s.cache.SetUsers(ctx, limit, filter, entries)
		}

		return entries, nil
	})
	if err != nil {
		return nil, err
	}

	entries, ok := v.([]model.LeaderboardEntry)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T", v)
	}
	return entries, nil
}

func (s *LeaderboardService) GetPlaceLeaderboard(ctx context.Context, limit int, filter model.LeaderboardFilter) ([]model.LeaderboardEntry, error) {
	if s.cache != nil {
		if entries, ok := s.cache.GetPlaces(ctx, limit, filter); ok {
			return entries, nil
		}
	}

	sfKey := fmt.Sprintf(
		"leaderboard:places:limit=%d:sort=%s:minr=%s:minrev=%d",
		limit,
		filter.SortBy,
		strconv.FormatFloat(filter.MinRating, 'f', -1, 64),
		filter.MinReviews,
	)

	v, err, _ := s.sf.Do(sfKey, func() (any, error) {
		if s.cache != nil {
			if entries, ok := s.cache.GetPlaces(ctx, limit, filter); ok {
				return entries, nil
			}
		}

		log.Println("DB QUERY leaderboard/places")

		entries, err := s.repo.GetTopPlaces(ctx, limit, filter)
		if err != nil {
			return nil, err
		}

		if s.cache != nil {
			s.cache.SetPlaces(ctx, limit, filter, entries)
		}

		return entries, nil

	})
	if err != nil {
		return nil, err
	}

	entries, ok := v.([]model.LeaderboardEntry)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T", v)
	}
	return entries, nil
}

func (s *LeaderboardService) GetBonusLeaderboard(ctx context.Context) ([]model.BonusLeaderboardEntry, error) {
	if s.cache != nil {
		if entries, ok := s.cache.GetBonuses(ctx); ok {
			return entries, nil
		}
	}

	const sfKey = "leaderboard:bonuses"

	v, err, _ := s.sf.Do(sfKey, func() (any, error) {
		if s.cache != nil {
			if entries, ok := s.cache.GetBonuses(ctx); ok {
				return entries, nil
			}
		}

		log.Println("DB QUERY leaderboard/bonuses")

		entries, err := s.repo.GetTopBonusUsers(ctx)
		if err != nil {
			return nil, err
		}

		if s.cache != nil {
			s.cache.SetBonuses(ctx, entries)
		}

		return entries, nil

	})
	if err != nil {
		return nil, err
	}

	entries, ok := v.([]model.BonusLeaderboardEntry)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T", v)
	}
	return entries, nil
}
