package main

import "time"

// +++ COUCHDB MONTHLY ARCHIVE +++

type monthDataArchive struct {
	ID  string `json:"_id"`
	Rev string `json:"_rev,omitempty"`

	Positions []position `json:"positions"`
}

type position struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Time      time.Time `json:"time"`
}

type positionByTime []position

func (b positionByTime) Len() int           { return len(b) }
func (b positionByTime) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b positionByTime) Less(i, j int) bool { return b[i].Time.Before(b[j].Time) }

// +++ COUCHDB USER FILE +++

type userDatabase struct {
	ID  string `json:"_id"`
	Rev string `json:"_rev,omitempty"`

	Users map[string]user `json:"users"`
}

type user struct {
	Name      string `json:"name"`
	EMail     string `json:"email"`
	Token     string `json:"token"`
	Protected bool   `json:"protected"`
	ViewToken string `json:"view_token"`
}

// +++ WEB DATA TRANSFER FORMAT FOR FRONTEND +++

type currentDataResponse struct {
	Now          currentDataTime `json:"now"`
	Date         currentDataTime `json:"date"`
	Distance     float64         `json:"distance"`
	TimeDelta    int64           `json:"timedelta"`
	Speed        float64         `json:"speed"`
	DisplaySpeed string          `json:"display_speed"`
	Latitude     float64         `json:"latitude"`
	Longitude    float64         `json:"longitude"`
	Place        string          `json:"place,omitempty"`
	Timestamp    int64           `json:"timestamp"`
}

type currentDataTime struct {
	time.Time
}

func (c currentDataTime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + c.Format("2006-01-02 15:04:05 MST") + "\""), nil
}
