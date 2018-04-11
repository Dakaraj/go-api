package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

// CustomServeMux is a struct that can be a multiplexer
type CustomServeMux struct {
}

// ServeHTTP is a function handler to be overridden
func (c *CustomServeMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		giveRandom(w, req)
		return
	}
	http.NotFound(w, req)
	return
}

func giveRandom(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Your random number is %f", rand.Float64())
}

func main() {
	mux := &CustomServeMux{}
	http.ListenAndServe(":8000", mux)
}
