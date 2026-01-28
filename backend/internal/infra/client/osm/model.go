package osm

type Place struct {
	Source   string // "osm"
	SourceID string // "way/73216328"
	Name     string
	Amenity  *string

	AddrCity        *string
	AddrCountry     *string
	AddrHouseNumber *string
	AddrPostcode    *string
	AddrStreet      *string

	Lat float64
	Lon float64
}

type SearchParams struct {
	City    string
	Search  *string
	Amenity *string
	Limit   int
	Offset  int // опционально
}
