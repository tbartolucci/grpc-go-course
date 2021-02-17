package main

import (
	"context"
	"fmt"
	"github.com/tbartolucci/udemy-grpc/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	fmt.Println("Hello I'm a Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial grpc server: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	fmt.Printf("Connected to client: %v\n", c)
	doUnary(c)
	doStream(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting : &greetpb.Greeting{
			FirstName: "Tom",
			LastName: "Bartolucci",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPF: %v", err)
	}

	log.Printf("Response: %v", res.Result)
}

func doStream(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting : &greetpb.Greeting{
			FirstName: "Bruzak",
			LastName: "Grinchy",
		},
	}
	resultStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPF: %v", err)
	}

	for {
		msg, err := resultStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}

		if err != nil {
			log.Fatalf("error while reading stream %v", err)
		}

		log.Printf("Response from stream: %v", msg.GetResult())
	}


}