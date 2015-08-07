package main

import "time"

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
}

type currentDataTime struct {
	time.Time
}

func (c currentDataTime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + c.Format("02.01.2006 15:04:05") + "\""), nil
}
