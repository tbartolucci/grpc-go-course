package main

import (
	"context"
	"fmt"
	"github.com/tbartolucci/udemy-grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) != 2 {
		panic("Please provide  two numbers to add!")
	}
	one,_ := strconv.Atoi(argsWithoutProg[0])
	two,_ := strconv.Atoi(argsWithoutProg[1])

	fmt.Printf("Let's find out what %d plus %d equals!\n", one, two)

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial grpc server: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	fmt.Printf("Connected to client: %v\n", c)
	doUnary(c, int64(one), int64(two))

}

func doUnary(c calculatorpb.CalculatorServiceClient, one int64, two int64) {
	req := &calculatorpb.SumRequest{
		AddendOne: one,
		AddendTwo: two,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}

	log.Printf("Response: %v", res.Result)
}

func doStream(c calculatorpb.CalculatorServiceClient, prime int64) {
	req :=  &calculatorpb.PrimeRequest{
		PrimeNumber: prime,
	}

	resultStream, err := c.Decompose(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Prime RPC: %v", err)
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