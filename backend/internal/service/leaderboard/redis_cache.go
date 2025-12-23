package leaderboard

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
	"github.com/redis/go-redis/v9"
)

type RedisLeaderboardCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewRedisLeaderboardCache(rdb *redis.Client, ttl time.Duration) *RedisLeaderboardCache {
	return &RedisLeaderboardCache{
		rdb: rdb,
		ttl: ttl,
	}
}

func (c *RedisLeaderboardCache) GetUsers(ctx context.Context, limit int, f model.LeaderboardFilter) ([]model.LeaderboardEntry, bool) {
	key := usersKey(limit, f)

	raw, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, false
	}
	if err != nil {
		return nil, false
	}
	if raw == "" {
		return nil, false
	}

	var entries []model.LeaderboardEntry
	if err := json.Unmarshal([]byte(raw), &entries); err != nil {
		return nil, false
	}

	return entries, true
}

func (c *RedisLeaderboardCache) SetUsers(ctx context.Context, limit int, f model.LeaderboardFilter, entries []model.LeaderboardEntry) {
	key := usersKey(limit, f)

	b, err := json.Marshal(entries)
	if err != nil {
		return
	}

	_ = c.rdb.Set(ctx, key, b, c.ttl).Err()
}

func (c *RedisLeaderboardCache) GetPlaces(ctx context.Context, limit int, f model.LeaderboardFilter) ([]model.LeaderboardEntry, bool) {
	key := placesKey(limit, f)

	raw, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil || raw == "" {
		return nil, false
	}
	if err != nil {
		return nil, false
	}

	var entries []model.LeaderboardEntry
	if err := json.Unmarshal([]byte(raw), &entries); err != nil {
		return nil, false
	}

	return entries, true
}

func (c *RedisLeaderboardCache) SetPlaces(ctx context.Context, limit int, f model.LeaderboardFilter, entries []model.LeaderboardEntry) {
	key := placesKey(limit, f)

	b, err := json.Marshal(entries)
	if err != nil {
		return
	}

	_ = c.rdb.Set(ctx, key, b, c.ttl).Err()
}

func (c *RedisLeaderboardCache) GetBonuses(ctx context.Context) ([]model.BonusLeaderboardEntry, bool) {
	key := bonusesKey()

	raw, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil || raw == "" {
		return nil, false
	}
	if err != nil {
		return nil, false
	}

	var entries []model.BonusLeaderboardEntry
	if err := json.Unmarshal([]byte(raw), &entries); err != nil {
		return nil, false
	}

	return entries, true
}

func (c *RedisLeaderboardCache) SetBonuses(ctx context.Context, entries []model.BonusLeaderboardEntry) {
	key := bonusesKey()

	b, err := json.Marshal(entries)
	if err != nil {
		return
	}

	_ = c.rdb.Set(ctx, key, b, c.ttl).Err()
}

func usersKey(limit int, f model.LeaderboardFilter) string {
	return fmt.Sprintf(
		"leaderboard:users:limit=%d:sort=%s:minr=%.2f:minrev=%d",
		limit, f.SortBy, f.MinRating, f.MinReviews,
	)
}

func placesKey(limit int, f model.LeaderboardFilter) string {
	return fmt.Sprintf(
		"leaderboard:places:limit=%d:sort=%s:minr=%.2f:minrev=%d",
		limit, f.SortBy, f.MinRating, f.MinReviews,
	)
}

func bonusesKey() string {
	return "leaderboard:bonuses"
}
