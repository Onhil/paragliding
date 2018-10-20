package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
)

var startTime time.Time

func main() {
	startTime = time.Now()
	Init()

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
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router)) // set listen port
}
