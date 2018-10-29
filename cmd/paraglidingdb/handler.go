package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Onhil/paragliding/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/p3lim/iso8601"
)

// MetaInfo for general meta information for project and api uptime
type MetaInfo struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

// Responds with current API staus
func getAPI(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, MetaInfo{
		Uptime:  iso8601.Format(time.Since(startTime)),
		Info:    "Service for paragliding tracks.",
		Version: "v1",
	})
}

// Returns all track IDs if any
func getTrackIDs(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, paragliding.GlobalDB.GetAll())
}

// Adds a new track to db
func postTrack(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
	} else if url, ok := data["url"]; !ok {
		http.Error(w, "Missing url", http.StatusBadRequest)
	} else if id, err := paragliding.GlobalDB.Add(url); err != nil {
		http.Error(w, "Url does not contain track data", http.StatusBadRequest)
	} else {
		response := make(map[string]int)
		response["id"] = id
		render.JSON(w, r, response)
	}
}

// Returns track with specific ID if existing
func getTrack(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if data, err := paragliding.GlobalDB.GetTrack(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		render.JSON(w, r, data)
	}
}

// Returns specific track field if ID and field exist
func getTrackField(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if data, err := paragliding.GlobalDB.GetTrack(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		field := chi.URLParam(r, "field")
		if fieldValue, err := data.GetField(field); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			render.JSON(w, r, fieldValue)
		}
	}
}

// Returns timestamp of latest added track
func getTickerLatest(w http.ResponseWriter, r *http.Request) {
	if track, err := paragliding.GlobalDB.TickerLatest(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		render.JSON(w, r, track.Timestamp)
	}
}

// Returns ticker info for tracks
func getTicker(w http.ResponseWriter, r *http.Request) {
	if ticker, err := paragliding.GlobalDB.Ticker(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		render.JSON(w, r, ticker)
	}
}

// Returns ticker infor for tracks after a timestamp
func getTickerTimestamp(w http.ResponseWriter, r *http.Request) {
	timestamp := chi.URLParam(r, "timestamp")
	if ticker, err := paragliding.GlobalDB.TickerTimestamp(timestamp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		render.JSON(w, r, ticker)
	}
}

// Adds a webhook to MongoDB with trigger info
func postWebhook(w http.ResponseWriter, r *http.Request) {
	var data paragliding.Webhooks
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
	}
	if data.MinTriggerValue == 0 {
		data.MinTriggerValue = 1
	}

	if data.WebhookURL == "" {
		http.Error(w, "Missing WebhookURL", http.StatusBadRequest)
	} else if id, err := paragliding.GlobalDB.AddWebhook(data.WebhookURL, data.MinTriggerValue); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		response := make(map[string]int)
		response["id"] = id
		render.JSON(w, r, response)
	}
}

// Returns webhook with id
func getWebhookID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "webhook_id")
	if webhook, err := paragliding.GlobalDB.GetWebhook(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		render.JSON(w, r, webhook)
	}
}

// Deletes webhook with id
func deleteWebhookID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "webhook_id")
	if webhook, err := paragliding.GlobalDB.DeleteWebhook(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		render.JSON(w, r, webhook)
	}
}

// Returns track collectiion count
func trackCount(w http.ResponseWriter, r *http.Request) {
	if count, err := paragliding.GlobalDB.TracksCount(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		response := make(map[string]int)
		response["count"] = count
		render.JSON(w, r, response)
	}
}

// Deletes track collection and returns count
func deleteAllTracks(w http.ResponseWriter, r *http.Request) {
	if count, err := paragliding.GlobalDB.DeleteAllTracks(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		response := make(map[string]int)
		response["count"] = count
		render.JSON(w, r, response)
	}
}
