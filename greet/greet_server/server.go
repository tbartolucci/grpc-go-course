package main

import (
	"context"
	"fmt"
	"github.com/tbartolucci/udemy-grpc/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
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

func (s *server) GreetManyTimes(request *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	log.Print("GreetManyTimes was invoked\n")
	firstName := request.Greeting.GetFirstName()
	for i := 0; i < 10;  i++ {
		result := fmt.Sprintf("Hello %s number %d\n", firstName, i)
		res := &greetpb.GreetManytimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (s *server) LongGreet(greetServer greetpb.GreetService_LongGreetServer) error {
	result := "Hello "
	for {
		req, err := greetServer.Recv()
		if err == io.EOF {
			//client stream is complete
			return greetServer.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += firstName + "! "
	}

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