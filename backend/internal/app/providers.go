package app

import (
	"github.com/redis/go-redis/v9"

	"github.com/kulikovroman08/reviewlink-backend/configs"
	clients "github.com/kulikovroman08/reviewlink-backend/internal/infra/client"
	newsclient "github.com/kulikovroman08/reviewlink-backend/internal/infra/client/news"
	"github.com/kulikovroman08/reviewlink-backend/internal/infra/client/osm"
	"github.com/kulikovroman08/reviewlink-backend/internal/infra/client/rss"
	svc "github.com/kulikovroman08/reviewlink-backend/internal/service"
)

func buildOSMClient(cfg *configs.Config, rdb *redis.Client) clients.OSMClient {
	base := osm.NewOverpassClient()
	return osm.NewCachedClient(
		base,
		rdb,
		cfg.OSMCacheTTL,
		cfg.OSMCacheHardTTL,
	)
}

func buildPublicNewsClient(cfg *configs.Config, rdb *redis.Client) svc.PublicNewsClient {
	rssURL := cfg.PublicNewsRSSURL
	if rssURL == "" {
		rssURL = "https://daily.afisha.ru/export/rss/google_newsstand/"
	}

	rssClient := rss.NewClient(rssURL)
	return newsclient.NewCachedClient(rssClient, rdb, cfg.PublicNewsCacheTTL)
}
