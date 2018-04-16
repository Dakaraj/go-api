package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer("/home/dakar/static")))
	// router.ServeFiles("/static/*filepath", http.Dir("/home/dakar/static"))
	log.Fatal(http.ListenAndServe(":80", router))
}
