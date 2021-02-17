package main

import (
	"context"
	"fmt"
	"github.com/tbartolucci/udemy-grpc/sum/sumpb"
	"google.golang.org/grpc"
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

	c := sumpb.NewSumServiceClient(cc)
	fmt.Printf("Connected to client: %v\n", c)
	req := &sumpb.SumRequest{
		AddendOne: int64(one),
		AddendTwo: int64(two),
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}

	log.Printf("Response: %v", res.Result)
}