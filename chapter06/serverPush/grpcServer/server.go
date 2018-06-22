package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc/reflection"

	pb "github.com/dakaraj/go-api/chapter06/serverPush/datafiles"
	"google.golang.org/grpc"
)

const (
	port      = ":50051"
	noOfSteps = 3
)

type server struct{}

func (s *server) MakeTransaction(in *pb.TransactionRequest, stream pb.MoneyTransaction_MakeTransactionServer) error {
	log.Println("Got request for money transfer...")
	log.Printf("Amount: $%f, From: A/c: %s, To: A/c: %s", in.Amount, in.From, in.To)
	// Send streams here
	for i := 0; i < noOfSteps; i++ {
		// Simulating I/O oe computation process using sleep...
		// Usually this will be saving money teransfer details in DB or talk to the third party API
		time.Sleep(time.Second * 2)
		// Once task is done, send the successful message back to the client
		if err := stream.Send(&pb.TransactionResponse{
			Status:      "good",
			Step:        int32(i),
			Description: fmt.Sprintf("Decryption of step %d", i),
		}); err != nil {
			log.Fatalf("%v.Send(%v) = %v", stream, "status", err.Error())
		}
	}
	log.Printf("Successfully transfered amount $%v from %v to %v", in.Amount, in.From, in.To)

	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Listening failed with error: ", err.Error())
	}
	// Create new gRPC server
	s := grpc.NewServer()
	// Register it with Proto service
	pb.RegisterMoneyTransactionServer(s, &server{})
	// Register reflection service on gRPC server
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to serve: ", err)
	}
}
