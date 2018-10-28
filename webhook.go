package main

import (
	"sort"

	"github.com/globalsign/mgo/bson"
)

// Webhooks stores data about a webhook
type Webhooks struct {
	ID              bson.ObjectId `bson:"_id,omitempty"`
	WebhookID       int           `json:"webhookid"`
	WebhookURL      string        `json:"WebhookURL"`
	MinTriggerValue int           `json:"MinTriggerValue"`
}

// WebhookIDs returns a slice wtih all WebhookID's
func WebhookIDs(db []Webhooks) []int {
	// Stores all existing ID's in a slice
	ids := make([]int, 0)
	for i := range db {
		ids = append(ids, db[i].WebhookID)
	}
	sort.Ints(ids)
	return ids
}

// CreateWebhook from url and mintriggervalue data
func CreateWebhook(url string, value int) (Webhooks, error) {
	var webhook Webhooks

	id := len(GlobalDB.GetAllWebhooks()) + 1
	webhook = Webhooks{
		WebhookID:       id,
		WebhookURL:      url,
		MinTriggerValue: value,
	}
	return webhook, nil
}
