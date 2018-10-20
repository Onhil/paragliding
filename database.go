package main

import (
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

	var all []int

	err = session.DB(db.DatabaseName).C(db.TrackCollectionName).Find(bson.M{}).All(&all)
	if err != nil {
		return []int{}
	}

	return all
}
