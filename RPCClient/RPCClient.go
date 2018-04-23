package main

import (
	"log"
	"net/rpc"
)

// Args struct for sending arguments to server
type Args struct{}

func main() {
	var reply int64
	args := Args{}
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Dialing error:", err)
	}
	err = client.Call("TimeServer.GiveServerTime", args, &reply)
	if err != nil {
		log.Fatal("Arith error:", err)
	}
	log.Printf("%d", reply)
}
