package main

import (
	"context"
	"fmt"
	"github.com/tbartolucci/udemy-grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type CalculatorServer struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (s *CalculatorServer) GetPrimes(n int64) []int64 {
	var list []int64
	k := int64(2)
	for n > 1 {
		if n % k == 0 { // if k evenly divides into N
			list = append(list, k)
			n = n / k // divide N by k so that we have the rest of the number left.
		} else {
			k = k + 1
		}
	}
	return list
}

func (s *CalculatorServer) Sum(c context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	response := &calculatorpb.SumResponse{
		Result: req.AddendOne + req.AddendTwo,
	}

	return response, nil
}

func (s *CalculatorServer) Decompose(request *calculatorpb.PrimeRequest,stream calculatorpb.CalculatorService_DecomposeServer) error {
	log.Print("Decompose was invoked\n")
	primes := s.GetPrimes(request.PrimeNumber)
	for i := 0; i < len(primes);  i++ {
		res := &calculatorpb.PrimeResponse{
			Result: primes[i],
		}
		stream.Send(res)
	}
	return nil
}

func (s * CalculatorServer) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	var sum int64
	var count int64
	sum = 0
	count = 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//client stream is complete
			return stream.SendAndClose(&calculatorpb.AverageResponse{
				Answer: float32(sum/count),
			})
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
		}
		sum = sum + req.GetNumber()
		count++
	}

	return nil
}


func main() {
	fmt.Println("Sum Server Starting....")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &CalculatorServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to CalculatorServer:  %v", err)
	}
}