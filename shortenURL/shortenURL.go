package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
)

type shortenURL struct {
	URL string
}

func filterContentType(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Currently in the content type checker middleware")
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("415 - Usupported Media Type. JSON required"))
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func createShortenedURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method Not Allowed"))
		return
	}
	var rURL shortenURL
	body := json.NewDecoder(r.Body)
	err := body.Decode(&rURL)
	if err != nil {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte("405 - Unsupported media type. Please send a valid JSON"))
	}
	_, pErr := url.ParseRequestURI(rURL.URL)
	if pErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request. Provided URL is not valid!"))
		return
	}
}

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.UseEncodedPath()

	r.HandleFunc("/api/v1/new", createShortenedURL)
	r.Queries("url")

	srv := http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:80",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
	}
	log.Fatal(srv.ListenAndServe())
}
