package main

import (
	"golang.org/x/net/context"
	"io"
	"log"

	pb "github.com/dakaraj/go-api/chapter06/serverPush/datafiles"
	"google.golang.org/grpc"
)

const address = "127.0.0.1:50051"

// ReceiveStream listens to stream contents and use them
func ReceiveStream(client pb.MoneyTransactionClient, req *pb.TransactionRequest) {
	log.Println("Started listening to server stream...")
	stream, err := client.MakeTransaction(context.Background(), req)
	if err != nil {
		log.Fatal("Could not get the stream data from server! Error: ", err.Error())
	}
	// Listen to stream of messages
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			// If there are no more messages
			break
		}
		if err != nil {
			log.Fatalf("%v.MakeTransaction(_) = _, %v", client, err.Error())
		}
		log.Printf("Status: %v, Operation: %v", response.Status, response.Description)
	}
}

func main() {
	// Set up a connection to server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Dialing to server failed with error: ", err.Error())
	}
	defer conn.Close()

	client := pb.NewMoneyTransactionClient(conn)

	// Prepare data. Get this from clients
	from := "1234"
	to := "5678"
	amount := float32(321.67)

	// Contact the server and print out its response
	ReceiveStream(client, &pb.TransactionRequest{
		From:   from,
		To:     to,
		Amount: amount,
	})
}
