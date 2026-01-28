package osm

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
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

	// 1) try cache
	if raw, err := c.rdb.Get(ctx, key).Bytes(); err == nil && len(raw) > 0 {
		var out []Place
		if err := json.Unmarshal(raw, &out); err == nil {
			return out, nil
		}
		// если кэш битый — просто игнорируем
	}

	// 2) fetch upstream
	out, err := c.inner.SearchPlaces(ctx, p)
	if err != nil {
		return nil, err
	}

	// 3) store cache (best-effort)
	if b, err := json.Marshal(out); err == nil {
		_ = c.rdb.Set(ctx, key, b, c.ttl).Err()
	}

	return out, nil
}

func (c *CachedClient) cacheKey(p SearchParams) string {
	// нормализуем, чтобы ключи были стабильными
	city := strings.TrimSpace(strings.ToLower(p.City))
	search := ""
	if p.Search != nil {
		search = strings.TrimSpace(strings.ToLower(*p.Search))
	}
	amenity := ""
	if p.Amenity != nil {
		amenity = strings.TrimSpace(strings.ToLower(*p.Amenity))
	}

	limit := p.Limit
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	base := fmt.Sprintf("city=%s|search=%s|amenity=%s|limit=%d", city, search, amenity, limit)

	sum := sha1.Sum([]byte(base))
	return "osm:places:" + hex.EncodeToString(sum[:])
}
