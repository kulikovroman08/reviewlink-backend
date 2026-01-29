package news

import (
	"context"
	"encoding/json"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/kulikovroman08/reviewlink-backend/internal/infra/client/rss"
	"github.com/redis/go-redis/v9"
)

const (
	defaultLimit = 3
	maxLimit     = 20

	cacheKey = "news:afisha:general:v1"
	source   = "afisha"
)

type CachedClient struct {
	inner *rss.Client
	rdb   *redis.Client
	ttl   time.Duration
}

func NewCachedClient(inner *rss.Client, rdb *redis.Client, ttl time.Duration) *CachedClient {
	return &CachedClient{
		inner: inner,
		rdb:   rdb,
		ttl:   ttl,
	}
}

type cachedPayload struct {
	CachedAt time.Time `json:"cached_at"`
	Items    []Item    `json:"items"`
}

func (c *CachedClient) ListNews(ctx context.Context, limit int) ([]Item, *time.Time, error) {
	if limit <= 0 {
		limit = defaultLimit
	}
	limit = clamp(limit, 1, maxLimit)

	// try cache
	if c.rdb != nil {
		if raw, err := c.rdb.Get(ctx, cacheKey).Bytes(); err == nil && len(raw) > 0 {
			var p cachedPayload
			if err := json.Unmarshal(raw, &p); err == nil {
				return take(p.Items, limit), &p.CachedAt, nil
			}
		}
	}

	// fetch fresh
	items, cachedAt, err := c.fetchAndMapAll(ctx)
	if err != nil {
		return nil, nil, err
	}

	// store cache (best-effort)
	if c.rdb != nil && cachedAt != nil {
		p := cachedPayload{CachedAt: *cachedAt, Items: items}
		if b, err := json.Marshal(p); err == nil {
			_ = c.rdb.Set(ctx, cacheKey, b, c.ttl).Err()
		}
	}

	return take(items, limit), cachedAt, nil
}

func (c *CachedClient) fetchAndMapAll(ctx context.Context) ([]Item, *time.Time, error) {
	raw, err := c.inner.Fetch(ctx)
	if err != nil {
		return nil, nil, err
	}
	log.Printf("[news] rss fetched items: %d", len(raw))

	out := make([]Item, 0)
	seen := make(map[string]struct{}, len(raw)) // dedupe by link

	for _, it := range raw {
		title := strings.TrimSpace(it.Title)
		link := strings.TrimSpace(it.Link)
		if title == "" || link == "" {
			continue
		}

		// ❌ выкидываем только тяжелое/негатив
		if matchAny(title, stopKeywords) {
			continue
		}

		tm, ok := parsePubDate(it.PubDate)
		if !ok {
			continue
		}

		if _, ok := seen[link]; ok {
			continue
		}
		seen[link] = struct{}{}

		out = append(out, Item{
			Title:       title,
			Link:        link,
			PublishedAt: tm,
			Source:      source,
		})
	}
	log.Printf("[news] rss passed after filters: %d", len(out))

	// newest first
	sort.Slice(out, func(i, j int) bool {
		return out[i].PublishedAt.After(out[j].PublishedAt)
	})

	now := time.Now()
	return out, &now, nil
}

func matchAny(s string, kws []string) bool {
	low := strings.ToLower(s)
	for _, k := range kws {
		if strings.Contains(low, strings.ToLower(k)) {
			return true
		}
	}
	return false
}

func parsePubDate(v string) (time.Time, bool) {
	v = strings.TrimSpace(v)
	if v == "" {
		return time.Time{}, false
	}

	// RSS классика
	if t, err := time.Parse(time.RFC1123Z, v); err == nil {
		return t, true
	}
	if t, err := time.Parse(time.RFC1123, v); err == nil {
		return t, true
	}

	// часто встречается у некоторых RSS
	if t, err := time.Parse(time.RFC3339, v); err == nil {
		return t, true
	}
	if t, err := time.Parse(time.RFC3339Nano, v); err == nil {
		return t, true
	}

	return time.Time{}, false
}

func take(items []Item, limit int) []Item {
	if len(items) <= limit {
		return items
	}
	return items[:limit]
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
