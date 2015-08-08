package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
)

var geocodingCache = map[string]osmNominatimResult{}

type osmNominatimResult struct {
	PlaceID     int64   `json:"place_id,string"`
	License     string  `json:"licence"`
	OSMType     string  `json:"osm_type"`
	OSMID       int64   `json:"osm_id,string"`
	Latitude    float64 `json:"lat,string"`
	Longitude   float64 `json:"lon,string"`
	DisplayName string  `json:"display_name"`
	Details     struct {
		HouseNo     string `json:"house_number"`
		Road        string `json:"road"`
		Suburb      string `json:"suburb"`
		Town        string `json:"town"`
		County      string `json:"county"`
		State       string `json:"state"`
		PostCode    string `json:"postcode"`
		Country     string `json:"country"`
		CountryCode string `json:"country_code"`
	} `json:"address"`
}

func (o osmNominatimResult) String() string {
	tplString := "{{ .DisplayName }}"

	switch o.Details.Country {
	case "Germany":
		tplString = "{{ .Details.Road }} {{ .Details.HouseNo }}, {{ .Details.PostCode }} {{ .Details.Town }}, {{ .Details.Country }}"
	}

	tpl, err := template.New("address").Parse(tplString)
	if err != nil {
		fmt.Printf("Unable to parse Template: %s\n", err)
	}
	buf := bytes.NewBuffer([]byte{})
	tpl.Execute(buf, o)

	return buf.String()
}

func geocodeCoordinate(lat, lon float64) string {
	cache := fmt.Sprintf("%.4f,%.4f", lat, lon)
	if cachedResult, ok := geocodingCache[cache]; ok {
		return cachedResult.String()
	}

	osmNominatim := fmt.Sprintf("http://nominatim.openstreetmap.org/reverse?format=json&lat=%.4f&lon=%.4f", lat, lon)

	req, _ := http.NewRequest("GET", osmNominatim, nil)
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("User-Agent", fmt.Sprintf("Mozilla/5.0 (compatible; Locationmaps/%s; +https://github.com/Luzifer/locationmaps)", version))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Unable to retrieve nominatim data: %s\n", err)
		return ""
	}
	defer resp.Body.Close()

	t := osmNominatimResult{}
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		fmt.Printf("Unable to parse nominatim data: %s\n", err)
		return ""
	}

	geocodingCache[cache] = t
	return t.String()
}
