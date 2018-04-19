package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func createShortenedURL(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
		return
	}
	query := r.URL.Query()
	fmt.Println(method)
	fmt.Println(query["url"][0])
}

func main() {
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.UseEncodedPath()

	router.HandleFunc("/api/v1/new", createShortenedURL)
	router.Queries("url")

	srv := http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:80",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
	}
	log.Fatal(srv.ListenAndServe())
}
