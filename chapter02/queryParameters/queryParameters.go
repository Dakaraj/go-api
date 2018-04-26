package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// QueryHandler function handles a request with query parameters
func QueryHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category name is: %s\n", query["category"][0])
	fmt.Fprintf(w, "Id is: %s\n", query["id"][0])
}

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/articles", QueryHandler)
	r.Queries("id", "category")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:80",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
	}
	log.Fatal(srv.ListenAndServe())
}
