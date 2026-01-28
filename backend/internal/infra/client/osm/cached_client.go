package osm

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type CachedClient struct {
	inner *OverpassClient
	rdb   *redis.Client
	ttl   time.Duration
}

func NewCachedClient(inner *OverpassClient, rdb *redis.Client, ttl time.Duration) *CachedClient {
	return &CachedClient{
		inner: inner,
		rdb:   rdb,
		ttl:   ttl,
	}
}

func (c *CachedClient) SearchPlaces(ctx context.Context, p SearchParams) ([]Place, error) {
	// если Redis отключен — просто ходим наружу
	if c.rdb == nil {
		return c.inner.SearchPlaces(ctx, p)
	}

	key := c.cacheKey(p)

	// 1) try exact cache
	if raw, err := c.rdb.Get(ctx, key).Bytes(); err == nil && len(raw) > 0 {
		var out []Place
		if err := json.Unmarshal(raw, &out); err == nil {
			return out, nil
		}
	}

	// ✅ если это поиск (search/amenity) — попробуем отфильтровать из кэша "каталога"
	hasSearch := p.Search != nil && strings.TrimSpace(*p.Search) != ""
	hasAmenity := p.Amenity != nil && strings.TrimSpace(*p.Amenity) != ""

	if hasSearch || hasAmenity {
		base := p
		base.Search = nil
		base.Amenity = nil

		baseKey := c.cacheKey(base)

		if raw, err := c.rdb.Get(ctx, baseKey).Bytes(); err == nil && len(raw) > 0 {
			var baseItems []Place
			if err := json.Unmarshal(raw, &baseItems); err == nil && len(baseItems) > 0 {
				filtered := filterLocalPlaces(baseItems, p)

				// best-effort: сохраним отфильтрованный результат по "точному" ключу
				if b, err := json.Marshal(filtered); err == nil {
					_ = c.rdb.Set(ctx, key, b, c.ttl).Err()
				}

				return filtered, nil
			}
		}
	}

	// 2) fetch upstream (как было)
	out, err := c.inner.SearchPlaces(ctx, p)
	if err != nil {
		return nil, err
	}

	// 3) store cache (best-effort) (как было)
	if b, err := json.Marshal(out); err == nil {
		_ = c.rdb.Set(ctx, key, b, c.ttl).Err()
	}

	return out, nil
}

func filterLocalPlaces(items []Place, p SearchParams) []Place {
	search := ""
	if p.Search != nil {
		search = strings.ToLower(strings.TrimSpace(*p.Search))
	}
	amenity := ""
	if p.Amenity != nil {
		amenity = strings.ToLower(strings.TrimSpace(*p.Amenity))
	}

	if search == "" && amenity == "" {
		return applyLimit(items, p.Limit)
	}

	out := make([]Place, 0, len(items))

	for _, it := range items {
		// search по имени (contains, без regex — устойчиво)
		if search != "" {
			name := strings.ToLower(strings.TrimSpace(it.Name))
			if !strings.Contains(name, search) {
				continue
			}
		}

		// amenity
		if amenity != "" {
			if it.Amenity == nil {
				continue
			}
			a := strings.ToLower(strings.TrimSpace(*it.Amenity))
			if a == "" || !strings.Contains(a, amenity) {
				continue
			}
		}

		out = append(out, it)
	}

	return applyLimit(out, p.Limit)
}

func applyLimit(items []Place, limit int) []Place {
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	if len(items) > limit {
		return items[:limit]
	}
	return items
}
