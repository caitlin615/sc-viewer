// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"encoding/json"
	"log"
	"time"
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) Broadcast(scribble *Scribble) {
	// TODO: Maybe store this somewhere and send last couple to new clients
	b, err := json.Marshal(scribble)
	if err != nil {
		log.Println(err)
		return
	}
	h.broadcast <- b
}

// TODO: Use Scribble in handwritingio/scribblechat/api/sc/metadata.go
type Scribble struct {
	ID        string            `json:"id"`
	DateSent  time.Time         `json:"date_sent"`
	ShareMode string            `json:"shareMode"`
	URLs      map[string]string `json:"urls"`

	Sent bool
}

// /webhook consumes from an AWS Lambda
// the Lambda function listens for mp4/gif object creations and send the data to the webhook
// which then adds broadcasts it
//
// from __future__ import print_function
//
// import json
// import urllib
// import boto3
// import os
// import logging
// from datetime import datetime
// import base64
//
// from urllib2 import Request, urlopen, URLError, HTTPError
//
// logging.info("Start logger")
// log = logging.getLogger()
// log.setLevel(logging.INFO)
//
// webhook_url = os.getenv("WEBHOOK_URL")
// secret = os.getenv("WEBHOOK_SECRET")
//
// def lambda_handler(event, context):
//     #print("Received event: " + json.dumps(event, indent=2))
//     items = []
//     for record in event['Records']:
//         print(record['eventID'])
//         if record['eventName'] != 'MODIFY':
//             continue
//        item = {}
//        date_sent = record['dynamodb']['NewImage'].get('date_sent')
//        if date_sent:
//            item["date_sent"] = date_sent
//         item["id"] = record['dynamodb']['NewImage']['id']['S']
//         item["shareMode"] = record['dynamodb']['NewImage']['exporterShareMode']['S']
//         urlRecord = record['dynamodb']['NewImage']['urls']['M']
//         item["urls"] = {}
//         for mimeType, value in urlRecord.iteritems():
//             item["urls"][mimeType] = value['S']
//
//         items.append(item)
//     try:
//         req = Request(webhook_url, json.dumps(items))
//         base64string = base64.b64encode('%s:%s' % ("lambda", secret))
//         req.add_header("Authorization", "Basic %s" % base64string)
//         try:
//             response = urlopen(req)
//             response.read()
//             # log.info("Success:", items)
//         except HTTPError as e:
//             log.error("Request failed: %d %s", e.code, e.reason)
//             raise e
//         except URLError as e:
//             log.error("Server connection failed: %s", e.reason)
//             raise e
//     except Exception as e:
//         log.info(e)
//         raise e
