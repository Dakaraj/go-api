package main

import (
	"io"
	"log"
	"net/http"
)

// Hello serves responce for "/hello" endpoint for application
func Hello(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, World!\n")
}

// Bye serves responce for "/bye" endpoint for application
func Bye(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Bye bye!\n")
}

func main() {
	http.HandleFunc("/hello", Hello)
	http.HandleFunc("/bye", Bye)
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}
