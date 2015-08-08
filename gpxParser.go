package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"time"
)

type gpx struct {
	XMLName xml.Name   `xml:"gpx"`
	Tracks  []gpxTrack `xml:"trk"`
}

type gpxTrack struct {
	XMLName  xml.Name     `xml:"trk"`
	Name     string       `xml:"name"`
	Segments []gpxSegment `xml:"trkseg"`
}

type gpxSegment struct {
	XMLName xml.Name        `xml:"trkseg"`
	Points  []gpxTrackPoint `xml:"trkpt"`
}

type gpxTrackPoint struct {
	XMLName   xml.Name  `xml:"trkpt"`
	Latitude  float64   `xml:"lat,attr"`
	Longitude float64   `xml:"lon,attr"`
	Elevation float64   `xml:"ele"`
	Time      time.Time `xml:"time"`
}

func parseGPXInput(source io.Reader) (*gpx, error) {
	out := gpx{}
	if err := xml.NewDecoder(source).Decode(&out); err != nil {
		return nil, err
	}

	return &out, nil
}

func extractPositionsFromGPX(source *gpx) []position {
	result := []position{}

	for _, track := range source.Tracks {
		for _, seg := range track.Segments {
			for _, pt := range seg.Points {
				result = append(result, position{
					Latitude:  pt.Latitude,
					Longitude: pt.Longitude,
					Time:      pt.Time,
				})
			}
		}
	}

	return result
}

func addPositionsToMonthArchives(user string, source []position) error {
	archives := map[string]*monthDataArchive{}

	for _, pos := range source {
		cacheString := fmt.Sprintf("%s_%d-%d", user, pos.Time.Year(), pos.Time.Month())
		if _, ok := archives[cacheString]; !ok {
			a, err := getMontlyArchive(user, pos.Time.Year(), int(pos.Time.Month()))
			if err != nil && !strings.Contains(err.Error(), "Code: 404") {
				return err
			}
			archives[cacheString] = a
		}

		archives[cacheString].Positions = append(archives[cacheString].Positions, pos)
	}

	for cs := range archives {
		saveMonthlyArchive(archives[cs])
	}

	return nil
}
