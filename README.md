### Paragliding

A web service to browse and upload .igc files


It's an RESTful API with 13 available calls:

## Track

- `GET /api`
	API status
- `GET /api/track`
	Array of all stored tracks
- `POST /api/track`
	Takes `{"url":"<url>"}` as JSON data and returns the assigned ID.
- `GET /api/igc/<id>`
	Returns track data (fields) for a valid `id`
- `GET /api/igc/<id>/<field>`
	Returns track `field` for a valid `id` and `field`

## Ticker
- `GET /api/ticker/latest`
    Latest timestamp
- `GET /api/ticker`
    Ticker info
- `GET /api/ticker/<timestamp>`
    Ticker info after timestamp

## Webhook
- `POST /api/webhook/new_track/`
    Takes `{"WebhookURL":     "<webhookurl>"`
          ` "MinTriggerValue": <number>`
          `}`
    and adds it to the database
- `GET /api/webhook/new_track/<id>`
    Returns webhook with that id
- `DELETE /api/webhook/new_track/<id>`
    Deletes webhook with that id


	
Demo of app
 - `https://paraglidingjg.herokuapp.com/`
