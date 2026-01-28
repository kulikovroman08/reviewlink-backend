package osm

type overpassResp struct {
	Elements []overpassElement `json:"elements"`
}

type overpassElement struct {
	Type   string            `json:"type"`
	ID     int64             `json:"id"`
	Lat    *float64          `json:"lat,omitempty"`
	Lon    *float64          `json:"lon,omitempty"`
	Center *overpassCenter   `json:"center,omitempty"`
	Tags   map[string]string `json:"tags,omitempty"`
}

type overpassCenter struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
