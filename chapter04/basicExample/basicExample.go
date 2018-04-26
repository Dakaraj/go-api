package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	gr "github.com/emicklei/go-restful"
)

func pingTime(req *gr.Request, res *gr.Response) {
	// Write time to the response
	io.WriteString(res, fmt.Sprintf("%s", time.Now()))
}

func main() {
	// Create a new web service
	webservice := new(gr.WebService)
	// Create a route and attach it to handler in the service
	webservice.Route(webservice.GET("/ping").To(pingTime))
	//Add the service to application
	gr.Add(webservice)
	http.ListenAndServe(":80", nil)
}
