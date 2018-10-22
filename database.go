package main

import (
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// GlobalDB interface for database use
var GlobalDB TrackStorage

// TrackMongoDB TODO
type TrackMongoDB struct {
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
		Key:        []string{"glider_id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = session.DB(db.DatabaseName).C(db.TrackCollectionName).EnsureIndex(index)
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

	err = session.DB(db.DatabaseName).C(db.TrackCollectionName).Find(bson.M{}).All(&tracks)
	if err != nil {
		return []int{}
	}

	return GetAll(tracks)
}

// Add TODO
func (db *TrackMongoDB) Add(url string) (int, error) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Parses url and returns track object
	track, err := Add(url)
	if err != nil {
		return 0, err
	}

	// Inserts track into mongoDB database
	err = session.DB(db.DatabaseName).C(db.TrackCollectionName).Insert(track)

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

	err = session.DB(db.DatabaseName).C(db.TrackCollectionName).Find(bson.M{"TrackID": id}).One(&track)
	if err != nil {
		return Track{}, err
	}

	return track, nil
}
