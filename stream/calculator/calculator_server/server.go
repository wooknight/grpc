package main

import (
	"context"
	"fmt"
	"grpc/stream/calculator/calculatorpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Received Sum RPC: %v", req)
	res := calculatorpb.SumResponse{
		SumResult: req.FirstNumber + req.SecondNumber,
	}
	return &res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	log.Printf("Received PrimeNumberDecomposition RPC : %v\n", req)
	num := req.GetNumber()
	divisor := int64(2)
	for num > 1 {
		if num%divisor == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			num = num / divisor
		} else {
			divisor++
			log.Printf("Divisor increased : %v\n", divisor)
		}
	}
	return nil
}

func main() {
	fmt.Println("Calculator server started")
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
