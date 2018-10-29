package paragliding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/globalsign/mgo"

	"github.com/globalsign/mgo/bson"
)

// Webhooks stores data about a webhook
type Webhooks struct {
	ID              bson.ObjectId `bson:"_id,omitempty"`
	WebhookID       int           `json:"webhookid"`
	WebhookURL      string        `json:"WebhookURL"`
	MinTriggerValue int           `json:"MinTriggerValue"`
	AddedSince      int           `json:"addedsince"`
	PrevTracksCount int           `json:"prevtrackscount"`
}

// CreateWebhook from url and mintriggervalue data
func CreateWebhook(url string, value int, c *mgo.Collection) (Webhooks, error) {
	var webhook Webhooks
	id, err := c.Count()
	if err != nil {
		return Webhooks{}, err
	}
	webhook = Webhooks{
		WebhookID:       id + 1,
		WebhookURL:      url,
		MinTriggerValue: value,
	}
	return webhook, nil
}

// SendMessage to a webhook about new track information
func SendMessage(webhook Webhooks, tc *mgo.Collection, wc *mgo.Collection) error {
	process := time.Now()
	var tracks []Track
	var newTracks []Track
	err := tc.Find(nil).All(&tracks)
	if err != nil {
		return err
	}
	// Appends new tracks added since last invoke
	for i := webhook.PrevTracksCount + 1; i < len(tracks); i++ {
		newTracks = append(newTracks, tracks[i])
	}
	ids := TrackIDs(newTracks)

	// Count of Track collection
	count, err := tc.Count()
	if err != nil {
		return err
	}
	// Creates message for webhook
	m := map[string]interface{}{
		"content": fmt.Sprintf("Latest timestamp: %s. New added tracks: %d [%d] (Processing time %s)",
			newTracks[len(newTracks)-1].Timestamp,
			count-webhook.PrevTracksCount, ids,
			time.Since(process)),
	}

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	// Sends message to webhook
	_, err = http.Post(webhook.WebhookURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	// Updates PrevTracksCount to current track collection count
	err = wc.Update(bson.M{}, bson.M{"$set": bson.M{"prevtrackscount": count}})
	if err != nil {
		return err
	}
	// Resets added since
	err = wc.Update(bson.M{}, bson.M{"$set": bson.M{"addedsince": 0}})
	return err
}
