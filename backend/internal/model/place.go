package model

import "github.com/google/uuid"

type PublicPlacesParams struct {
	City    string
	Search  *string
	Amenity *string
	Limit   int
	Offset  int
}

type PublicPlacesResult struct {
	Items []PublicPlace
	Total *int
}

type PublicPlace struct {
	Source   string
	SourceID string
	Name     string
	Amenity  *string

	Address  PublicPlaceAddress
	Location PublicLatLon

	Reviewlink *PublicPlaceReviewlinkMeta
}

type PublicPlaceAddress struct {
	City        *string
	Street      *string
	HouseNumber *string
	Postcode    *string
	Country     *string
	Display     *string
}

type PublicLatLon struct {
	Lat float64
	Lon float64
}

type PublicPlaceReviewlinkMeta struct {
	PlaceID      uuid.UUID
	Rating       *float64
	ReviewsCount int
	HasOwner     bool
}
