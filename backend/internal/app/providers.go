package app

import (
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/kulikovroman08/reviewlink-backend/configs"
	clients "github.com/kulikovroman08/reviewlink-backend/internal/infra/client"
	newsclient "github.com/kulikovroman08/reviewlink-backend/internal/infra/client/news"
	"github.com/kulikovroman08/reviewlink-backend/internal/infra/client/osm"
	"github.com/kulikovroman08/reviewlink-backend/internal/infra/client/rss"
)

func buildOSMClient(cfg *configs.Config, rdb *redis.Client) clients.OSMClient {
	base := osm.NewOverpassClient()
	return osm.NewCachedClient(base, rdb, 60*time.Minute)
}

func buildPublicNewsClient(cfg *configs.Config, rdb *redis.Client) clients.PublicNewsClient {
	rssURL := "https://daily.afisha.ru/export/rss/google_newsstand/"

	rssClient := rss.NewClient(rssURL)
	return newsclient.NewCachedClient(rssClient, rdb, 60*time.Second)
}
