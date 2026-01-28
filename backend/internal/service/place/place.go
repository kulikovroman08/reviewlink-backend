package place

import (
	"context"
	"fmt"
	"log"
	"strings"

	clients "github.com/kulikovroman08/reviewlink-backend/internal/infra/client"
	osm "github.com/kulikovroman08/reviewlink-backend/internal/infra/client/osm"

	"github.com/google/uuid"
	"github.com/kulikovroman08/reviewlink-backend/configs"
	"github.com/kulikovroman08/reviewlink-backend/internal/repository"
	"github.com/kulikovroman08/reviewlink-backend/internal/service/token"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

type placeService struct {
	placeRepo    repository.PlaceRepository
	tokenService *token.Service
	osmClient    clients.OSMClient
	cfg          *configs.Config
}

func NewPlaceService(
	placeRepo repository.PlaceRepository,
	tokenService *token.Service,
	osmClient clients.OSMClient,
	cfg *configs.Config,
) *placeService {
	return &placeService{
		placeRepo:    placeRepo,
		tokenService: tokenService,
		osmClient:    osmClient,
		cfg:          cfg,
	}
}

func (s *placeService) CreatePlace(ctx context.Context, place model.Place) (*model.Place, error) {
	if place.OwnerID == nil {
		return nil, fmt.Errorf("owner_id is required")
	}

	if place.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if place.Address == "" {
		return nil, fmt.Errorf("address is required")
	}

	place.ID = uuid.New()

	if err := s.placeRepo.CreatePlace(ctx, &place); err != nil {
		return nil, fmt.Errorf("failed to create place: %w", err)
	}

	count := s.cfg.TokensAutoCount
	if _, err := s.tokenService.GenerateTokens(ctx, place.OwnerID.String(), place.ID.String(), count); err != nil {
		fmt.Printf("failed to auto-generate tokens for place %s: %v\n", place.ID, err)
	}

	return &place, nil
}

func (s *placeService) GetPlacesByOwner(ctx context.Context, ownerID string) ([]model.Place, error) {
	return s.placeRepo.GetPlacesByOwner(ctx, ownerID)
}

func (s *placeService) ListPublicPlaces(ctx context.Context, params model.PublicPlacesParams) (model.PublicPlacesResult, error) {
	if params.City == "" {
		return model.PublicPlacesResult{}, fmt.Errorf("city is required")
	}

	limit := params.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	osmPlaces, err := s.osmClient.SearchPlaces(ctx, osm.SearchParams{
		City:    params.City,
		Search:  params.Search,
		Amenity: params.Amenity,
		Limit:   limit,
		Offset:  params.Offset,
	})
	if err != nil {
		log.Printf("OSM SearchPlaces error: %+v", err)
		return model.PublicPlacesResult{}, fmt.Errorf("osm search: %w", err)
	}

	if params.Offset > 0 && params.Offset < len(osmPlaces) {
		osmPlaces = osmPlaces[params.Offset:]
	}
	if len(osmPlaces) > limit {
		osmPlaces = osmPlaces[:limit]
	}

	sourceIDs := make([]string, 0, len(osmPlaces))
	for _, p := range osmPlaces {
		sourceIDs = append(sourceIDs, p.SourceID)
	}

	metaMap, err := s.placeRepo.GetPublicMetaBySourceIDs(ctx, "osm", sourceIDs)
	if err != nil {
		return model.PublicPlacesResult{}, fmt.Errorf("db meta: %w", err)
	}

	items := make([]model.PublicPlace, 0, len(osmPlaces))
	for _, p := range osmPlaces {
		it := model.PublicPlace{
			Source:   p.Source,
			SourceID: p.SourceID,
			Name:     p.Name,
			Amenity:  p.Amenity,
			Address: model.PublicPlaceAddress{
				City:        p.AddrCity,
				Street:      p.AddrStreet,
				HouseNumber: p.AddrHouseNumber,
				Postcode:    p.AddrPostcode,
				Country:     p.AddrCountry,
			},
			Location: model.PublicLatLon{Lat: p.Lat, Lon: p.Lon},
		}

		it.Address.Display = buildDisplayAddress(it.Address)

		if m, ok := metaMap[p.SourceID]; ok {
			it.Reviewlink = &model.PublicPlaceReviewlinkMeta{
				PlaceID:      m.PlaceID,
				Rating:       m.Rating,
				ReviewsCount: m.ReviewsCount,
				HasOwner:     m.HasOwner,
			}
		}

		items = append(items, it)
	}

	return model.PublicPlacesResult{
		Items: items,
		Total: nil,
	}, nil
}

func buildDisplayAddress(a model.PublicPlaceAddress) *string {
	parts := make([]string, 0, 3)

	if a.City != nil && *a.City != "" {
		parts = append(parts, *a.City)
	}

	line := ""
	if a.Street != nil && *a.Street != "" {
		line = *a.Street
	}
	if a.HouseNumber != nil && *a.HouseNumber != "" {
		if line != "" {
			line = line + ", " + *a.HouseNumber
		} else {
			line = *a.HouseNumber
		}
	}
	if line != "" {
		parts = append(parts, line)
	}

	if len(parts) == 0 {
		return nil
	}
	s := strings.Join(parts, " — ")
	return &s
}

func (s *placeService) EnsureFromPublic(
	ctx context.Context,
	source string,
	sourceID string,
	name string,
) (string, error) {

	if source == "" || sourceID == "" {
		return "", fmt.Errorf("source and sourceID are required")
	}

	if name == "" {
		return "", fmt.Errorf("name is required")
	}

	return s.placeRepo.EnsureFromPublic(ctx, source, sourceID, name)
}
