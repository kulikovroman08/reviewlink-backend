package osm

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type CachedClient struct {
	inner   *OverpassClient
	rdb     *redis.Client
	softTTL time.Duration
	hardTTL time.Duration
}

func NewCachedClient(inner *OverpassClient, rdb *redis.Client, softTTL, hardTTL time.Duration) *CachedClient {
	return &CachedClient{
		inner:   inner,
		rdb:     rdb,
		softTTL: softTTL,
		hardTTL: hardTTL,
	}
}

type cachedPayload struct {
	FetchedAt time.Time `json:"fetched_at"`
	Items     []Place   `json:"items"`
}

func (c *CachedClient) SearchPlaces(ctx context.Context, p SearchParams) ([]Place, error) {
	if c.rdb == nil {
		return c.inner.SearchPlaces(ctx, p)
	}

	key := c.cacheKey(p)
	now := time.Now()

	if cp, ok := c.getPayload(ctx, key); ok && len(cp.Items) > 0 {
		if now.Sub(cp.FetchedAt) > c.softTTL {
			go c.refreshInBackground(key, p)
		}
		return applyLimit(cp.Items, p.Limit), nil
	}

	hasSearch := p.Search != nil && strings.TrimSpace(*p.Search) != ""
	hasAmenity := p.Amenity != nil && strings.TrimSpace(*p.Amenity) != ""

	if hasSearch || hasAmenity {
		base := p
		base.Search = nil
		base.Amenity = nil

		baseKey := c.cacheKey(base)

		if baseCP, ok := c.getPayload(ctx, baseKey); ok && len(baseCP.Items) > 0 {
			if now.Sub(baseCP.FetchedAt) > c.softTTL {
				go c.refreshInBackground(baseKey, base)
			}

			filtered := filterLocalPlaces(baseCP.Items, p)

			_ = c.setPayload(ctx, key, cachedPayload{FetchedAt: time.Now(), Items: filtered})

			return filtered, nil
		}
	}

	pFetch := p
	pFetch.Limit = maxLimit
	out, err := c.inner.SearchPlaces(ctx, pFetch)
	if err != nil {
		return nil, err
	}

	_ = c.setPayload(ctx, key, cachedPayload{FetchedAt: time.Now(), Items: out})

	return applyLimit(out, p.Limit), nil
}

func (c *CachedClient) refreshInBackground(key string, p SearchParams) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	lockKey := key + ":lock"

	ok, err := c.rdb.SetNX(ctx, lockKey, "1", 3*time.Minute).Result()
	if err != nil || !ok {
		return
	}
	defer func() { _ = c.rdb.Del(ctx, lockKey).Err() }()

	if cp, ok := c.getPayload(ctx, key); ok && time.Since(cp.FetchedAt) <= c.softTTL {
		return
	}

	pFetch := p
	pFetch.Limit = maxLimit
	items, err := c.inner.SearchPlaces(ctx, pFetch)
	if err != nil {
		return
	}

	_ = c.setPayload(ctx, key, cachedPayload{FetchedAt: time.Now(), Items: items})
}

func (c *CachedClient) getPayload(ctx context.Context, key string) (cachedPayload, bool) {
	raw, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil || len(raw) == 0 {
		return cachedPayload{}, false
	}

	var cp cachedPayload
	if err := json.Unmarshal(raw, &cp); err != nil {
		return cachedPayload{}, false
	}
	return cp, true
}

func (c *CachedClient) setPayload(ctx context.Context, key string, cp cachedPayload) error {
	b, err := json.Marshal(cp)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, b, c.hardTTL).Err()
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
		if search != "" {
			name := strings.ToLower(strings.TrimSpace(it.Name))
			if !strings.Contains(name, search) {
				continue
			}
		}

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
