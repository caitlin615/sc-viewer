package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/caitlin615/sc-viewer/ws"
)

var T = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		layoutData := struct{}{}
		err := T.ExecuteTemplate(w, "index.html", layoutData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	hub := ws.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		ws.Serve(hub, w, r)
	})

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		_, secret, ok := r.BasicAuth()
		if !ok || secret != os.Getenv("WEBHOOK_SECRET") {
			w.Header().Set("WWW-Authenticate", `Basic realm="Enter the admin secret as password"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		var scribble ws.Scribble
		if err := json.Unmarshal(body, &scribble); err != nil {
			log.Println(err, string(body))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		log.Printf("scribble: %#+v\n", scribble)

		hub.Broadcast(&scribble)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8100"
	}
	address := "0.0.0.0:" + port
	log.Println("now listening on", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
