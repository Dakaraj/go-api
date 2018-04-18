package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// ArticleHandler is a route handler
func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	// mux.Vars returns all path variables as a map
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category is: %v\n", vars["category"])
	fmt.Fprintf(w, "Article Id is: %v\n", vars["id"])
}

// QueryHandler is a route handler
func QueryHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Got parameter category: %v\n", query["category"][0])
	fmt.Fprintf(w, "Got parameter id: %v\n", query["id"][0])
}

func main() {
	// Create a new router
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.UseEncodedPath()
	router.Path("/articles/{category}/{id:[0-9]+}").HandlerFunc(ArticleHandler)

	router.HandleFunc("/articles", QueryHandler)
	router.Queries("category", "id")

	server := http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:80",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
	}
	log.Fatal(server.ListenAndServe())
}
