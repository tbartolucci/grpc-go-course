package main

import (
	"context"
	"fmt"
	"github.com/tbartolucci/udemy-grpc/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{
	// Embed the unimplemented server
	greetpb.UnimplementedGreetServiceServer
}

func (s *server) Greet(ctx context.Context, request *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function invoked\n")
	firstName := request.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (s *server) GreetManyTimes(request *greetpb.GreetManyTimesRequest, timesServer greetpb.GreetService_GreetManyTimesServer) error {
	panic("implement me")
}

func (s *server) LongGreet(greetServer greetpb.GreetService_LongGreetServer) error {
	panic("implement me")
}

func (s *server) GreetEveryone(everyoneServer greetpb.GreetService_GreetEveryoneServer) error {
	panic("implement me")
}

func (s *server) GreetWithDeadline(ctx context.Context, request *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	panic("implement me")
}

func main(){
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server:  %v", err)
	}

}