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

// DatabaseInterface to better manage database functions
type DatabaseInterface interface {
	Init()
	GetAll() []int
	Add(url string) (int, error)
	GetTrack(idURL string) (Track, error)
	TickerLatest() (Track, error)
	Ticker() (Ticker, error)
	TickerTimestamp(ts string) (Ticker, error)
	GetAllWebhooks() []int
	AddWebhook(url string, value int) (int, error)
	GetWebhook(id string) (Webhooks, error)
	DeleteWebhook(id string) (Webhooks, error)
	TracksCount() (int, error)
	DeleteAllTracks() (int, error)
}

// GlobalDB interface for database use
var GlobalDB DatabaseInterface

// TrackMongoDB is a struct with all neccessary MongoDB info
type TrackMongoDB struct {
	DatabaseInfo          mgo.DialInfo
	DatabaseURL           string
	DatabaseName          string
	TrackCollectionName   string
	WebhookCollectionName string
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

	// TODO make sure you cannot add the same url twice

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

// GetAllWebhooks gets all
func (db *TrackMongoDB) GetAllWebhooks() []int {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var webhooks []Webhooks

	collection := session.DB(db.DatabaseName).C(db.WebhookCollectionName)

	err = collection.Find(nil).All(&webhooks)
	if err != nil {
		return []int{}
	}
	return WebhookIDs(webhooks)
}

// AddWebhook inserts webhook information into MongoDB
func (db *TrackMongoDB) AddWebhook(url string, value int) (int, error) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// TODO make sure you cannot add the same url twice

	// Creates webhook from data
	webhook, err := CreateWebhook(url, value)
	if err != nil {
		return 0, err
	}

	// Inserts webhook into mongoDB database
	err = session.DB(db.DatabaseName).C(db.WebhookCollectionName).
		Insert(webhook)
	if err != nil {
		return 0, err
	}

	return webhook.WebhookID, nil
}

// GetWebhook returns webhook with a specific id
func (db *TrackMongoDB) GetWebhook(id string) (Webhooks, error) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var webhook Webhooks

	ids, err := strconv.Atoi(id)
	if err != nil {
		return Webhooks{}, errors.New("Invalid ID")
	}

	// Returns webhook with specific id
	collection := session.DB(db.DatabaseName).C(db.WebhookCollectionName)
	err = collection.Find(bson.M{"webhookid": ids}).One(&webhook)
	if err != nil {
		return Webhooks{}, err
	}

	return webhook, nil
}

// DeleteWebhook removes webhook from MongoDB with a specific id
func (db *TrackMongoDB) DeleteWebhook(id string) (Webhooks, error) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	webhook, err := db.GetWebhook(id)

	// Removes webhook with a specific id
	collection := session.DB(db.DatabaseName).C(db.WebhookCollectionName)
	err = collection.Remove(bson.M{"webhookid": webhook.WebhookID})
	if err != nil {
		return Webhooks{}, err
	}

	return webhook, nil
}

// TracksCount returns count of Track collection
func (db *TrackMongoDB) TracksCount() (int, error) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Returns length of Tracks
	ids, err := session.DB(db.DatabaseName).C(db.TrackCollectionName).Count()
	if err != nil {
		return 0, err
	}

	return ids, nil
}

// DeleteAllTracks WARNING!!!
// Deletes all Tracks from MongoDB
func (db *TrackMongoDB) DeleteAllTracks() (int, error) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	collection := session.DB(db.DatabaseName).C(db.TrackCollectionName)

	// Returns count of Track collection
	ids, err := collection.Count()
	if err != nil {
		return 0, err
	}
	// Removes all Tracks from MongoDB
	collection.RemoveAll(nil)
	return ids, nil
}
