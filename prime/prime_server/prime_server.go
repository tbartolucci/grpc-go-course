package main

import (
	"fmt"
	"github.com/tbartolucci/udemy-grpc/prime/primepb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	primepb.UnimplementedPrimeServiceServer
}

func (s *server) GetPrimes(n int64) []int64 {
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

func (s *server) Decompose(request *primepb.PrimeRequest,stream primepb.PrimeService_DecomposeServer) error {
	log.Print("Decompose was invoked\n")
	primes := s.GetPrimes(request.PrimeNumber)
	for i := 0; i < len(primes);  i++ {
		res := &primepb.PrimeResponse{
			Result: primes[i],
		}
		stream.Send(res)
	}
	return nil
}


func main() {
	fmt.Println("Prime Server Starting....")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	primepb.RegisterPrimeServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server:  %v", err)
	}
}