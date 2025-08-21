package place

import (
	"context"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

type PlaceService interface {
	CreatePlace(ctx context.Context, place model.Place) (*model.Place, error)
}
