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
	if len(argsWithoutProg) == 0 {
		panic("Usage: program [method] [options...]")
	}

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial grpc server: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	fmt.Printf("Connected to client: %v\n", c)

	if argsWithoutProg[0] == "sum" {
		if len(argsWithoutProg) != 3 {
			panic("You must provide two numbers for summing!")
		}
		one,_ := strconv.Atoi(argsWithoutProg[1])
		two,_ := strconv.Atoi(argsWithoutProg[2])
		fmt.Printf("Let's find out what %d plus %d equals!\n", one, two)
		doUnary(c, int64(one), int64(two))

	} else if argsWithoutProg[0] == "prime" {
		if len(argsWithoutProg) != 2 {
			panic("You must provide a single number for prime decomposition!")
		}
		num,_ := strconv.Atoi(argsWithoutProg[1])
		fmt.Printf("Let's find out what the prime decomposition of %d is!\n", num)
		doStream(c, int64(num))
	}
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