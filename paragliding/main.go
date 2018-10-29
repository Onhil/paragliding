package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Onhil/paragliding"
	"github.com/go-chi/chi"
)

var startTime time.Time

func main() {
	startTime = time.Now()
	paragliding.GlobalDB = &paragliding.TrackMongoDB{
		DatabaseURL:           "mongodb://admin:admin1@ds145562.mlab.com:45562/paragliding",
		DatabaseName:          "paragliding",
		TrackCollectionName:   "Tracks",
		WebhookCollectionName: "Webhooks",
	}
	paragliding.Paging = 5
	paragliding.GlobalDB.Init()
	defer paragliding.Sessions.Close()
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
				r.Get("/latest/", getTickerLatest)
				r.Get("/", getTicker)
				r.Get("/{timestamp:[A-Za-z0-9_]+}/", getTickerTimestamp)
			})
			r.Route("/webhook", func(r chi.Router) {
				r.Route("/new_track", func(r chi.Router) {
					r.Post("/", postWebhook)
					r.Route("/{webhook_id:[0-9]+}", func(r chi.Router) {
						r.Get("/", getWebhookID)
						r.Delete("/", deleteWebhookID)
					})
				})
			})
		})
	})
	//
	router.Route("/admin", func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Get("/tracks_count/", trackCount)
			r.Delete("/tracks/", deleteAllTracks)
		})
	})
	//log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router)) // set listen port

	// localtesting
	log.Fatal(http.ListenAndServe(":8080", router)) // set listen port
}
