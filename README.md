sc-viewer
-----------

Real-time displaying of videos or gifs via a websocket connection.

Remote URLs are consumed by `/webhook` as json
```json
[
  {
    "id": "someID",
    "date_sent": "someDateSent",
    "shareMode": "yourShareMode",
    "urls": {
      "gif": "https://mycdn/gif.gif",
      "mp4": "https://mycdn/mp4.mp4"
    }
  }
  ...
]
```

The URLs are then broadcast to all clients connected via websocket.

## Note:
There is currently no persistence. Once the URL is pushed into the websocket, it is forgotten.
Clients that refresh the page will start from scratch. IDs are stored in memory for deduping.

It is only configured to keep 30 videos/gifs on the DOM at a time.

# Setup

Requires go 1.7 and above

    go get ./...
    go run main.go

# Variables
* `PORT`: webserver port (default: `8100`)
* `WEBHOOK_SECRET`: Basic authentication for `/webhook` endpoint (ignores username)

# Examples

See [examples/dynamodb.py](examples/dynamodb.py) for an example of how to send
data from an AWS DynamoDB database (via an AWS Lambda function) to the `/websocket`
endpoint.
