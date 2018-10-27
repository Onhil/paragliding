package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/globalsign/mgo"
	"github.com/go-chi/chi"
)

var startTime time.Time
var session *mgo.Session

func main() {
	startTime = time.Now()
	mongoDialInfo := &mgo.DialInfo{
		Addrs:    []string{os.Getenv("mongodb://@ds145562.mlab.com:45562/paragliding")},
		Timeout:  60 * time.Second,
		Database: os.Getenv("paragliding"),
		Username: os.Getenv("admin"),
		Password: os.Getenv("admin1"),
	}
	GlobalDB = &TrackMongoDB{
		*mongoDialInfo,
		"mongodb://admin:admin1@ds145562.mlab.com:45562/paragliding",
		"paragliding",
		"Tracks",
	}
	GlobalDB.Init()
	router := chi.NewRouter()
	router.Route("/paragliding", func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Get("/", getAPI)
			r.Route("/track", func(r chi.Router) {
				r.Get("/", getTrackIDs)
				r.Post("/", postTrack)
				r.Route("/{id:[0-9]+}", func(r chi.Router) {
					r.Get("/", getTrack)
					r.Get("/{field:[A-Za-z_]+}/", getTrackField)
				})
			})
			r.Route("/ticker", func(r chi.Router) {
				//r.Get("/latest/", )
				//r.get("/{timestamp:[0-9]+}/", )
			})
			r.Route("/weebhook", func(r chi.Router) {
				r.Route("/new_track", func(r chi.Router) {
					//r.Post("/", )
					r.Route("/webhook_id:[0-9]+}", func(r chi.Router) {
						//r.Get("/", )
						//r.Delete("/", )
					})
				})
			})
		})
	})
	//
	router.Route("/admin", func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			//r.Get("/tracks_count/", )
			//r.Delete("/tracks/", )
		})
	})
	//log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router)) // set listen port

	// localtesting
	log.Fatal(http.ListenAndServe(":8080", router)) // set listen port
}
