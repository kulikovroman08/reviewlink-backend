package osm

import "fmt"

func buildOverpassQuery(p SearchParams, limit int) string {
	amenity := `restaurant|cafe`
	if p.Amenity != nil && *p.Amenity != "" {
		amenity = *p.Amenity
	}

	nameFilter := ""
	if p.Search != nil && *p.Search != "" {
		nameFilter = fmt.Sprintf(`["name"~"%s",i]`, *p.Search)
	}

	city := p.City // для читаемости

	return fmt.Sprintf(`
[out:json][timeout:25];

// 1) Ищем административную границу как relation
(
  rel["boundary"="administrative"]["name"="%[1]s"];
  rel["boundary"="administrative"]["name:ru"="%[1]s"];
  rel["boundary"="administrative"]["name"~"%[1]s",i];
  rel["boundary"="administrative"]["name:ru"~"%[1]s",i];
)->.r;

// 2) relation -> area (ВАЖНО: синтаксис именно такой)
.r map_to_area -> .a;

// 3) Ищем POI в области
(
  node["amenity"~"%[2]s"]%[3]s(area.a);
  way["amenity"~"%[2]s"]%[3]s(area.a);
  relation["amenity"~"%[2]s"]%[3]s(area.a);
);

out center %[4]d;
`, city, amenity, nameFilter, limit)
}
