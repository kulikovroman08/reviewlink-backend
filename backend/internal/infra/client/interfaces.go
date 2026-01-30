package clients

import (
	"context"

	"github.com/kulikovroman08/reviewlink-backend/internal/infra/client/osm"
	"github.com/kulikovroman08/reviewlink-backend/internal/infra/client/rss"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock_clients.go -package=mocks
type OSMClient interface {
	SearchPlaces(ctx context.Context, p osm.SearchParams) ([]osm.Place, error)
}

type NewsClient interface {
	Fetch(ctx context.Context) ([]rss.Item, error)
}
