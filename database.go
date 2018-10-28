package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// GlobalDB interface for database use
var GlobalDB TrackStorage

// TrackMongoDB is a struct with all neccessary MongoDB info
type TrackMongoDB struct {
	DatabaseInfo        mgo.DialInfo
	DatabaseURL         string
	DatabaseName        string
	TrackCollectionName string
}

// Init intializes MongoDB
func (db *TrackMongoDB) Init() {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	index := mgo.Index{
		Key:        []string{"TrackID"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = session.DB(db.DatabaseName).C(db.TrackCollectionName).
		EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

// GetAll makes a slice with all TrackID's
func (db *TrackMongoDB) GetAll() []int {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var tracks []Track

	// Puts all tracks in a Track slice
	err = session.DB(db.DatabaseName).C(db.TrackCollectionName).
		Find(bson.M{}).All(&tracks)
	if err != nil {
		return []int{}
	}

	// Returns a slice with all TrackID's
	return TrackIDs(tracks)
}

// Add a track to the mongoDB
func (db *TrackMongoDB) Add(url string) (int, error) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Parses url and returns track object
	track, err := Parse(url)
	if err != nil {
		return 0, err
	}

	// Inserts track into mongoDB database
	err = session.DB(db.DatabaseName).C(db.TrackCollectionName).
		Insert(track)

	if err != nil {
		fmt.Printf("error in Insert(): %v", err.Error())
		return 0, err
	}

	return track.TrackID, nil
}

// GetTrack gets track wid a a specific TrackID
func (db *TrackMongoDB) GetTrack(id string) (Track, error) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Converts the id string to int
	var track Track
	ids, err := strconv.Atoi(id)
	if err != nil {
		return track, errors.New("Invalid ID")
	}
	// Finds track with a TrackID
	err = session.DB(db.DatabaseName).C(db.TrackCollectionName).
		Find(bson.M{"trackid": ids}).One(&track)
	if err != nil {
		return Track{}, err
	}
	return track, nil
}

//TickerLatest returns latest added track
func (db *TrackMongoDB) TickerLatest() (Track, error) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var track Track

	// Returns the latest added track
	collection := session.DB(db.DatabaseName).C(db.TrackCollectionName)
	if size, err := collection.Count(); err != nil {
		log.Fatal("Error: ", err)
	} else {
		err = collection.Find(nil).Skip(size - 1).One(&track)
		if err != nil {
			return Track{}, err
		}
	}

	return track, nil
}

// Ticker returns the timestamp of the first added track and the first and last
// timestamp for the "paging". Lastly it returns the proccessing time
func (db *TrackMongoDB) Ticker() (Ticker, error) {
	proccess := time.Now()
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var latest Track
	var start Track
	stop := make([]Track, paging, paging)

	collection := session.DB(db.DatabaseName).C(db.TrackCollectionName)
	if size, err := collection.Count(); err != nil {
		log.Fatal("Error: ", err)
	} else {
		// Finds the latest added track
		err = collection.Find(nil).Skip(size - 1).One(&latest)
		if err != nil {
			return Ticker{}, err
		}
		// Finds the first added track
		err = collection.Find(nil).One(&start)
		if err != nil {
			return Ticker{}, err
		}
		// Makes a Track slice with length paging
		err = collection.Find(nil).SetMaxScan(paging).All(&stop)
		if err != nil {
			return Ticker{}, err
		}
	}

	ticker := Ticker{
		Latest:     latest.Timestamp,
		Start:      start.Timestamp,
		Stop:       stop[len(stop)-1].Timestamp,
		Tracks:     TrackIDs(stop),
		Processing: time.Since(proccess),
	}
	return ticker, nil
}

// TickerTimestamp TODO
func (db *TrackMongoDB) TickerTimestamp(ts string) (Ticker, error) {
	proccess := time.Now()
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var latest Track
	var start Track
	stop := make([]Track, paging, paging)

	collection := session.DB(db.DatabaseName).C(db.TrackCollectionName)
	if size, err := collection.Count(); err != nil {
		log.Fatal("Error: ", err)
	} else {
		// Finds the latest added track
		err = collection.Find(nil).Skip(size - 1).One(&latest)
		if err != nil {
			return Ticker{}, err
		}
		// Finds track with a specific timestamp
		err = collection.Find(bson.M{"timestamp": bson.ObjectIdHex(ts)}).One(&start)
		if err != nil {
			return Ticker{}, err
		}
		// Makes a Track slice with length paging from Track start
		err = collection.Find(nil).SetMaxScan(start.TrackID + paging).All(&stop)
		if err != nil {
			return Ticker{}, err
		}
	}
	ticker := Ticker{
		Latest:     latest.Timestamp,
		Start:      stop[0].Timestamp,
		Stop:       stop[len(stop)-1].Timestamp,
		Tracks:     TrackIDs(stop),
		Processing: time.Since(proccess),
	}
	return ticker, nil
}
