package main

import (
	"log"
	"net"

	pb "github.com/dakaraj/go-api/chapter06/grpc_example/datafiles"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const port = ":50051"

type server struct{}

// MakeTransaction implements MoneyTransactionServer.MakeTransaction
func (s *server) MakeTransaction(ctx context.Context, in *pb.TransactionRequest) (*pb.TransactionResponse, error) {
	log.Println("Got request for money transfer....")
	log.Printf("Amount: %f, From A/c:%s, To A/c:%s", in.Amount, in.From, in.To)

	// DB logic goes here
	return &pb.TransactionResponse{Confirmation: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Failed to listen with error: ", err.Error())
	}
	s := grpc.NewServer()
	pb.RegisterMoneyTransactionServer(s, &server{})
	// Register reflection service on gRPC server
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
