package main

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/marni/goigc"
)

// TrackStorage TODO
type TrackStorage interface {
	Init()
	GetAll() []int
	Add(url string) (int, error)
	GetTrack(idURL string) (Track, error)
}

// Track data
type Track struct {
	ID             bson.ObjectId `bson:"_id,omitempty"`
	TrackID        int           `json:"trackid"`
	HDate          time.Time     `json:"H_date"`
	Pilot          string        `json:"pilot"`
	Glider         string        `json:"glider"`
	GliderID       string        `json:"glider_id"`
	TrackLength    float64       `json:"track_length"`
	TrackSourceURL string        `json:"track_src_url"`
}

// Ticker data
type Ticker struct {
	Latest     string        `json:"t_latest"`
	Start      string        `json:"t_start"`
	Stop       string        `json:"t_stop"`
	Tracks     []Track       `json:"tracks"`
	Processing time.Duration `json:"processing"`
}

// TrackIDs TODO
func TrackIDs(db []Track) []int {
	// Stores all existing ID's in a slice
	ids := make([]int, 0)
	for i := range db {
		ids = append(ids, db[i].TrackID)
	}
	sort.Ints(ids)
	return ids
}

// Parse TODO
func Parse(url string) (Track, error) {
	track, err := igc.ParseLocation(url)
	if err != nil {
		return Track{}, err
	}

	// Calulates total distance of track data
	distance := 0.0
	for i := 0; i < len(track.Points)-1; i++ {
		distance += track.Points[i].Distance(track.Points[i+1])
	}
	id := len(GlobalDB.GetAll()) + 1

	trac := Track{
		TrackID:        id,
		HDate:          track.Date,
		Pilot:          track.Pilot,
		Glider:         track.GliderType,
		GliderID:       track.GliderID,
		TrackLength:    distance,
		TrackSourceURL: url,
	}
	return trac, nil
}

/*
// GetTrack TODO
func GetTrack(id string) (Track, error) {
	var track Track

	// Converts ID to int
	id, err := strconv.Atoi(id)
	if err != nil {
		return track, errors.New("Invalid ID")
	}

	// Returns track data if ID exists
	data, ok := db[id]
	if ok {
		return data, nil
	}
	return track, errors.New("Track ID " + id + " does not exist")
}
*/

// GetField TODO
func (track *Track) GetField(field string) (string, error) {
	// Returns the field that matches the given struct json tag
	value := reflect.ValueOf(track).Elem()
	for i := 0; i < value.NumField(); i++ {
		if value.Type().Field(i).Tag.Get("json") == field {
			return fmt.Sprint(value.Field(i)), nil
		}

	}
	return "", errors.New("Track has no field " + field)
}
