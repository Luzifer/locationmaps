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

	resp := CurrentDataResponse{
		Now:       CurrentDataTime{Time: time.Now().UTC()},
		Date:      CurrentDataTime{Time: current.Time.UTC()},
		Distance:  haversine(current.Longitude, current.Latitude, last.Longitude, last.Latitude),
		TimeDelta: current.Time.Unix() - last.Time.Unix(),
		Latitude:  current.Latitude,
		Longitude: current.Longitude,
		Place:     geocodeCoordinate(current.Latitude, current.Longitude),
		Timestamp: current.Time.UTC().Unix(),
	}
	resp.Speed = resp.Distance / (float64(resp.TimeDelta) / 3600.0)
	resp.DisplaySpeed = fmt.Sprintf("%.2f km/h", resp.Speed)

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(resp); err != nil {
		fmt.Printf("Unable to encode response: %s\n", err)
	}
}

func handleGPXAdd(res http.ResponseWriter, r *http.Request) {
	usr, ok := users.Users[r.FormValue("user")]
	if !ok {
		http.Error(res, "User not found.", http.StatusNotFound)
		return
	}

	if r.FormValue("token") != usr.Token {
		http.Error(res, "Invalid token.", http.StatusForbidden)
		return
	}

	file, _, err := r.FormFile("gpxfile")
	if err != nil {
		http.Error(res, "Unable to read GPX file", http.StatusBadRequest)
		return
	}

	g, err := parseGPXInput(file)
	if err != nil {
		http.Error(res, fmt.Sprintf("Unable to parse GPX file: %s", err), http.StatusInternalServerError)
		return
	}
	err = addPositionsToMonthArchives(usr.Name, extractPositionsFromGPX(g))
	if err != nil {
		http.Error(res, fmt.Sprintf("Unable to save GPX positions: %s", err), http.StatusInternalServerError)
		return
	}

	http.Error(res, "OK", http.StatusOK)
}

func handleSimpleAdd(res http.ResponseWriter, r *http.Request) {
	usr, ok := users.Users[r.FormValue("user")]
	if !ok {
		http.Error(res, "User not found.", http.StatusNotFound)
		return
	}

	if r.FormValue("token") != usr.Token {
		http.Error(res, "Invalid token.", http.StatusForbidden)
		return
	}

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

func parseCoordinate(coord string) float64 {
	r, err := strconv.ParseFloat(coord, 64)
	if err != nil {
		return 0.0
	}

	return r
}
