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

func main() {
	// Create a new riuter
	router := mux.NewRouter()
	router.Path("/articles/{category}/{id:[0-9]+}").HandlerFunc(ArticleHandler)
	server := http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:80",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
	}
	log.Fatal(server.ListenAndServe())
}
