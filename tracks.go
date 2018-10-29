package paragliding

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/globalsign/mgo"

	"github.com/globalsign/mgo/bson"
	"github.com/marni/goigc"
)

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
	Timestamp      bson.ObjectId `bson:"timestamp"`
}

// Ticker data
type Ticker struct {
	Latest     bson.ObjectId `json:"t_latest"`
	Start      bson.ObjectId `json:"t_start"`
	Stop       bson.ObjectId `json:"t_stop"`
	Tracks     []int         `json:"tracks"`
	Processing time.Duration `json:"processing"`
}

// TrackIDs makes a slice with all TrackID's
func TrackIDs(db []Track) []int {
	// Stores all existing ID's in a slice
	ids := make([]int, 0)
	for i := range db {
		ids = append(ids, db[i].TrackID)
	}
	sort.Ints(ids)
	return ids
}

// Parse url and makes a Track
func Parse(url string, c *mgo.Collection) (Track, error) {
	track, err := igc.ParseLocation(url)
	if err != nil {
		return Track{}, err
	}

	// Calulates total distance of track data
	distance := 0.0
	for i := 0; i < len(track.Points)-1; i++ {
		distance += track.Points[i].Distance(track.Points[i+1])
	}

	// Returns count of Track Collection
	id, err := c.Count()
	if err != nil {
		return Track{}, err
	}

	trac := Track{
		TrackID:        id + 1,
		HDate:          track.Date,
		Pilot:          track.Pilot,
		Glider:         track.GliderType,
		GliderID:       track.GliderID,
		TrackLength:    distance,
		TrackSourceURL: url,
		Timestamp:      bson.NewObjectIdWithTime(time.Now()),
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

// GetField gets a specific json field in a Track
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
