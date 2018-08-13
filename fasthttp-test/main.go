package main

import (
	"fmt"
	"time"

	f "github.com/valyala/fasthttp"
)

var wcfmcClient = f.PipelineClient{
	Addr:         "www.whocanfixmycar.com",
	IsTLS:        true,
	ReadTimeout:  10 * time.Second,
	WriteTimeout: 10 * time.Second,
}

// var wcfmcClient = f.Client{
// 	ReadTimeout:  10 * time.Second,
// 	WriteTimeout: 10 * time.Second,
// }

func main() {
	fmt.Println("Starting an HTTP request")
	req := f.AcquireRequest()
	res := f.AcquireResponse()

	req.SetRequestURI("https://www.whocanfixmycar.com")
	req.Header.SetMethod("GET")
	req.Header.Add("Accept", "application/json")

	wcfmcClient.Do(req, res)

	sCode := res.StatusCode()
	bodyBytes := res.Body()
	headers := res.Header.String()

	// sCode, body, _ := client.Get(nil, "https://whocanfixmycar.com")
	fmt.Printf("Code: %d\nheaders: %s\nbody: %s\n", sCode, headers, string(bodyBytes))
}
