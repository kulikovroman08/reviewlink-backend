package osm

import (
	"fmt"
)

func mapElementToPlace(el overpassElement) (Place, bool) {
	lat, lon, ok := extractLatLon(el)
	if !ok {
		return Place{}, false
	}

	sourceID := fmt.Sprintf("%s/%d", el.Type, el.ID)

	name := ""
	amenity := ""
	if el.Tags != nil {
		name = el.Tags["name"]
		amenity = el.Tags["amenity"]
	}
	if name == "" {
		name = fallbackName(amenity, el.Tags)
	}

	var amenityPtr *string
	if amenity != "" {
		a := amenity
		amenityPtr = &a
	}

	pl := Place{
		Source:   "osm",
		SourceID: sourceID,
		Name:     name,
		Amenity:  amenityPtr,
		Lat:      lat,
		Lon:      lon,
	}

	if el.Tags != nil {
		pl.AddrCity = strPtrOrNil(el.Tags["addr:city"])
		pl.AddrCountry = strPtrOrNil(el.Tags["addr:country"])
		pl.AddrHouseNumber = strPtrOrNil(el.Tags["addr:housenumber"])
		pl.AddrPostcode = strPtrOrNil(el.Tags["addr:postcode"])
		pl.AddrStreet = strPtrOrNil(el.Tags["addr:street"])
	}

	return pl, true
}

func extractLatLon(el overpassElement) (float64, float64, bool) {
	if el.Lat != nil && el.Lon != nil {
		return *el.Lat, *el.Lon, true
	}
	if el.Center != nil {
		return el.Center.Lat, el.Center.Lon, true
	}
	return 0, 0, false
}

func strPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	ss := s
	return &ss
}

func fallbackName(amenity string, tags map[string]string) string {
	street := ""
	house := ""
	if tags != nil {
		street = tags["addr:street"]
		house = tags["addr:housenumber"]
	}
	if street != "" && house != "" {
		if amenity != "" {
			return fmt.Sprintf("%s, %s %s", amenity, street, house)
		}
		return fmt.Sprintf("%s %s", street, house)
	}
	if amenity != "" {
		return amenity
	}
	return "place"
}
