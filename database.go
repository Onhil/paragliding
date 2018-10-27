package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// GlobalDB interface for database use
var GlobalDB TrackStorage

// TrackMongoDB TODO
type TrackMongoDB struct {
	DatabaseInfo        mgo.DialInfo
	DatabaseURL         string
	DatabaseName        string
	TrackCollectionName string
}

// Init TODO
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

// GetAll TODO
func (db *TrackMongoDB) GetAll() []int {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var tracks []Track

	err = session.DB(db.DatabaseName).C(db.TrackCollectionName).
		Find(bson.M{}).All(&tracks)
	if err != nil {
		return []int{}
	}

	return TrackIDs(tracks)
}

// Add TODO
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

// GetTrack TODO
func (db *TrackMongoDB) GetTrack(id string) (Track, error) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var track Track
	ids, err := strconv.Atoi(id)
	if err != nil {
		return track, errors.New("Invalid ID")
	}
	err = session.DB(db.DatabaseName).C(db.TrackCollectionName).
		Find(bson.M{"trackid": ids}).One(&track)
	if err != nil {
		return Track{}, err
	}
	return track, nil
}

// ObjectID(_id).getTimestamp()

//TickerLatest TODO
func (db *TrackMongoDB) TickerLatest() (Track, error) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var track Track

	err = session.DB(db.DatabaseName).C(db.TrackCollectionName).
		Find(nil).Skip(len(GlobalDB.GetAll()) - 1).One(&track)
	if err != nil {
		return Track{}, err
	}
	return track, nil
}
