sc-viewer
-----------

Real-time displaying of videos or gifs via a websocket connection.

Remote URLs are consumed by `/webhook` as json
```json
{
  "last_modified": "date",
  "url": "https://mycdn/gif.gif",
}
```

The URL is then broadcast into all clients connected via websocket.

## Note:
There is currently no persistence. Once the URL is pushed into the websocket, it is forgotten.
Clients that refresh the page will start from scratch.

It is only configured to keep 30 videos/gifs on the DOM at a time.

# Setup

Requires go 1.7 and above

    go get ./...
    go run main.go

# Variables
* `PORT`: webserver port
* `WEBHOOK_SECRET`: Basic authentication for `/webhook` endpoint (ignores username)
