package place

import (
	"context"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

type PlaceRepository interface {
	CreatePlace(ctx context.Context, place *model.Place) error
}
