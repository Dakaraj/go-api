package main

import (
	"fmt"
	"time"

	f "github.com/valyala/fasthttp"
)

var client = f.PipelineClient{
	// Name:         "BandStack/Golang/HTTPClient/email:kramarev.anton@gmail.com",
	ReadTimeout:  10 * time.Second,
	WriteTimeout: 10 * time.Second,
}

func main() {
	fmt.Println("Starting an HTTP request")
	req := f.AcquireRequest()
	req.SetRequestURI("https://whocanfixmycar.com")
	req.Header.Add("Custom", "Header")
	client.Do(req, resp)

	// sCode, body, _ := client.Get(nil, "https://whocanfixmycar.com")
	fmt.Printf("Code: %d\nbody: %s\n", sCode, string(body[:100]))
}
