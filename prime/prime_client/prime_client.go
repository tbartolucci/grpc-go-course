package main

import (
	"context"
	"fmt"
	"github.com/tbartolucci/udemy-grpc/prime/primepb"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) != 1 {
		panic("Please provide  two numbers to add!")
	}
	num,_ := strconv.Atoi(argsWithoutProg[0])

	fmt.Printf("Let's find out what primes are in %d!\n", num)

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial grpc server: %v", err)
	}
	defer cc.Close()

	c := primepb.NewPrimeServiceClient(cc)
	fmt.Printf("Connected to client: %v\n", c)
	req := &primepb.PrimeRequest{
		PrimeNumber: int64(num),
	}

	res, err := c.Decompose(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}

	for {
		msg, err := res.Recv()
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