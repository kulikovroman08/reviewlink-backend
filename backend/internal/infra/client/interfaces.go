package clients

import (
	"context"

	"github.com/kulikovroman08/reviewlink-backend/internal/infra/client/osm"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock_clients.go -package=mocks
type OSMClient interface {
	SearchPlaces(ctx context.Context, p osm.SearchParams) ([]osm.Place, error)
}
