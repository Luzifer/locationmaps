package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func handleGetLatest(res http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	archive, err := getMontlyArchive(
		vars["user"],
		time.Now().Year(),
		int(time.Now().Month()),
	)

	if err != nil {
		fmt.Printf("Unable to load month archive: %s\n", err)
		http.Error(res, "Unable to load month archive", http.StatusInternalServerError)
		return
	}

	if len(archive.Positions) == 0 {
		http.Error(res, "Found no positions", http.StatusNotFound)
		return
	}

	sort.Sort(sort.Reverse(positionByTime(archive.Positions)))

	var current, last position

	current = archive.Positions[0]
	if len(archive.Positions) > 1 {
		last = archive.Positions[1]
	} else {
		last = archive.Positions[0]
	}

	resp := currentDataResponse{
		Now:       currentDataTime{Time: time.Now()},
		Date:      currentDataTime{Time: current.Time},
		Distance:  haversine(current.Longitude, current.Latitude, last.Longitude, last.Latitude),
		TimeDelta: current.Time.Unix() - last.Time.Unix(),
		Latitude:  current.Latitude,
		Longitude: current.Longitude,
	}
	resp.Speed = resp.Distance / (float64(resp.TimeDelta) / 3600.0)
	resp.DisplaySpeed = fmt.Sprintf("%.2f km/h", resp.Speed)

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(resp); err != nil {
		fmt.Printf("Unable to encode response: %s\n", err)
	}
}

func handleSimpleAdd(res http.ResponseWriter, r *http.Request) {
	ll := strings.Split(r.FormValue("location"), ",")
	if len(ll) != 2 {
		http.Error(res, "Invalid location", http.StatusBadRequest)
		return
	}

	entry := position{
		Time:      time.Now(),
		Latitude:  parseCoordinate(ll[0]),
		Longitude: parseCoordinate(ll[1]),
	}

	if entry.Latitude == 0.0 || entry.Longitude == 0.0 {
		http.Error(res, fmt.Sprintf("Invalid coordinates: %.6f,%.6f", entry.Latitude, entry.Longitude), http.StatusBadRequest)
		return
	}

	archive, err := getMontlyArchive(
		r.FormValue("user"),
		time.Now().Year(),
		int(time.Now().Month()),
	)

	if err != nil && !strings.Contains(err.Error(), "Code: 404") {
		fmt.Printf("Unable to load month archive: %s\n", err)
		http.Error(res, "Unable to load month archive", http.StatusInternalServerError)
		return
	}

	archive.Positions = append(archive.Positions, entry)

	err = saveMonthlyArchive(archive)

	if err != nil {
		fmt.Printf("Unable to save month archive: %s\n", err)
		http.Error(res, "Unable to save month archive", http.StatusInternalServerError)
		return
	}

	http.Error(res, "OK", http.StatusOK)
}

func getMontlyArchive(user string, year, month int) (*monthDataArchive, error) {
	doc := &monthDataArchive{
		ID:        fmt.Sprintf("%s_%d-%d", user, year, month),
		Positions: []position{},
	}

	if err := couchConn.Get(doc.ID, doc); err != nil {
		return doc, err
	}

	return doc, nil
}

func saveMonthlyArchive(archive *monthDataArchive) error {
	res, err := couchConn.Save(archive)

	if err != nil {
		return err
	}

	if !res.Ok {
		return fmt.Errorf("Could not store archive: %s / %s", res.Reason, res.Error)
	}

	archive.Rev = res.Rev

	return nil
}

func parseCoordinate(coord string) float64 {
	r, err := strconv.ParseFloat(coord, 64)
	if err != nil {
		return 0.0
	}

	return r
}
