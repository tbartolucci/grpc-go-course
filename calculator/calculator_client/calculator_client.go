package main

import (
	"context"
	"fmt"
	"github.com/tbartolucci/udemy-grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		panic("Usage: program [method] [options...]")
	}

	certFile := "ssl/ca.crt"
	creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
	if sslErr != nil {
		log.Fatalf("Error loading crt file: %v", sslErr)
	}
	opts := grpc.WithTransportCredentials(creds)
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Failed to dial grpc server: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	fmt.Printf("Connected to client: %v\n", c)

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if argsWithoutProg[0] == "sum" {
		if len(argsWithoutProg) != 3 {
			panic("You must provide two numbers for summing!")
		}
		one,_ := strconv.Atoi(argsWithoutProg[1])
		two,_ := strconv.Atoi(argsWithoutProg[2])
		fmt.Printf("Let's find out what %d plus %d equals!\n", one, two)
		doUnary(c, ctx, int64(one), int64(two))

	} else if argsWithoutProg[0] == "prime" {
		if len(argsWithoutProg) != 2 {
			panic("You must provide a single number for prime decomposition!")
		}
		num,_ := strconv.Atoi(argsWithoutProg[1])
		fmt.Printf("Let's find out what the prime decomposition of %d is!\n", num)
		handleStream(c, ctx, int64(num))
	} else if argsWithoutProg[0] == "average" {
		doStream(c, ctx, argsWithoutProg[1:])
	} else if argsWithoutProg[0] == "max" {
		doBiDi(c, ctx, argsWithoutProg[1:])
	} else if argsWithoutProg[0] == "sqrt" {
		num,_ := strconv.Atoi(argsWithoutProg[1])
		doErrorUnary(c, ctx, int64(num))
	}
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient, ctx context.Context, number int64) {
	res, err := c.SquareRoot(ctx, &calculatorpb.SquareRootRequest{ Number: number})
	if err != nil {
		respErr,ok := status.FromError(err)
		if ok {
			// actual error from GRPC, user error
			if respErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline exceeded.")
			} else if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number.")
			}
			fmt.Printf("[%v] %s\n", respErr.Code(), respErr.Message())
			return
		} else {
			// non standard error
			log.Fatalf("Big error calling squareroot: %v", err)
			return
		}
	}
	fmt.Printf("Result of Square root of %v: %v\n", number, res.GetNumber())
}

func doBiDi(c calculatorpb.CalculatorServiceClient, ctx context.Context, numbers []string) {
	fmt.Println("Do Bidirectional Streaming")
	// create a stream by invoking the cient
	stream, err := c.Maximum(ctx)
	if err != nil {
		log.Fatalf("error while creating stream: %v", err)
		return
	}

	waitc := make(chan struct{})
	// we send a bunch of messages to the client (go routine)
	go func() {
		for _, numString := range numbers{
			num,_ := strconv.Atoi(numString)
			req := &calculatorpb.MaximumRequest{ Number : int64(num) }
			fmt.Printf("Sending number: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// we receive a bunch of message the client (go routine)
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving: %v", err)
				break
			}
			fmt.Printf("Current Max: %v\n", res.GetResult())
		}
		close(waitc)
	}()
	// block until everything is done
	<-waitc
}
// Send Add Request
func doUnary(c calculatorpb.CalculatorServiceClient, ctx context.Context, one int64, two int64) {
	req := &calculatorpb.SumRequest{
		AddendOne: one,
		AddendTwo: two,
	}
	res, err := c.Sum(ctx, req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}

	log.Printf("Response: %v", res.Result)
}
// Send Prime Request
func handleStream(c calculatorpb.CalculatorServiceClient, ctx context.Context, prime int64) {
	req :=  &calculatorpb.PrimeRequest{
		PrimeNumber: prime,
	}

	resultStream, err := c.Decompose(ctx, req)
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
// Send Average Request
func doStream(c calculatorpb.CalculatorServiceClient, ctx context.Context, args []string) {
	fmt.Println("Starting to stream to the server....")
	stream, err := c.Average(ctx)
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	var requests []*calculatorpb.AverageRequest
	for _,numString := range args {
		num,_ := strconv.Atoi(numString)
		requests = append(requests, &calculatorpb.AverageRequest{
			Number : int64(num),
		})
	}

	for _, req := range requests {
		log.Printf("Sending %v request\n", req.Number )
		stream.Send(req)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Printf("Response: %v\n", res.Answer)
}