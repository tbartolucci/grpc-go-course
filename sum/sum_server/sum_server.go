package main

import (
	"context"
	"fmt"
	"github.com/tbartolucci/udemy-grpc/sum/sumpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	sumpb.UnimplementedSumServiceServer
}

func (s *server) Sum(c context.Context, req *sumpb.SumRequest) (*sumpb.SumResponse, error) {
	response := &sumpb.SumResponse{
		Result: req.AddendOne + req.AddendTwo,
	}

	return response, nil
}

func main() {
	fmt.Println("Sum Server Starting....")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	sumpb.RegisterSumServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server:  %v", err)
	}
}