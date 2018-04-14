package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir("/home/dakar/static"))
	log.Fatal(http.ListenAndServe(":80", router))
}
