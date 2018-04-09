package main

import (
	"io"
	"log"
	"net/http"
)

// MyServer handles "/hello" endpoint for application
func MyServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, World!\n")
}

func ByeFunc(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Bye bye!")
}

func main() {
	http.HandleFunc("/hello", MyServer)
	http.HandleFunc("/bye", ByeFunc)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
