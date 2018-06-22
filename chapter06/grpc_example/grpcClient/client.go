package main

import (
	"context"
	"log"

	pb "github.com/dakaraj/go-api/chapter06/grpc_example/datafiles"
	"google.golang.org/grpc"
)

const address = "127.0.0.1:50051"

func main() {
	// Set up a connection to server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Dialing to server failed with error:", err.Error())
	}
	defer conn.Close()

	c := pb.NewMoneyTransactionClient(conn)

	// Prepare data. Get this from clients on Frontend
	from := "1234"
	to := "5678"
	amount := float32(321.67)

	// Contact the server and print out its response
	resp, err := c.MakeTransaction(context.Background(), &pb.TransactionRequest{
		From:   from,
		To:     to,
		Amount: amount,
	})
	if err != nil {
		log.Fatal("Transaction failed with error: ", err.Error())
	}
	log.Println("Transaction confirmed with status: ", resp.Confirmation)
}
