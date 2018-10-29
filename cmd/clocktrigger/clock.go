package main

import (
	"log"
	"time"

	"github.com/Onhil/paragliding/database"
	"github.com/globalsign/mgo"
)

// MongoDB TODO
type MongoDB struct {
	DatabaseURL           string
	DatabaseName          string
	TrackCollectionName   string
	WebhookCollectionName string
}

func main() {
	var webhooks []paragliding.Webhooks
	GlobalDB := &MongoDB{
		DatabaseURL:           "mongodb://admin:admin1@ds145562.mlab.com:45562/paragliding",
		DatabaseName:          "paragliding",
		TrackCollectionName:   "Tracks",
		WebhookCollectionName: "Webhooks",
	}
	session, err := mgo.Dial(GlobalDB.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	wc := session.DB(GlobalDB.DatabaseName).C(GlobalDB.WebhookCollectionName)
	tc := session.DB(GlobalDB.DatabaseName).C(GlobalDB.TrackCollectionName)

	err = wc.Find(nil).All(&webhooks)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for i := range webhooks {
			if webhooks[i].AddedSince >= webhooks[i].MinTriggerValue {
				err = paragliding.SendMessage(webhooks[i], wc, tc)
				if err != nil {
					log.Fatal(err)
				}
			}

		}
		time.Sleep(10 * time.Minute)
	}

}
