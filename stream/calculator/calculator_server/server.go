package main

import (
	"context"
	"fmt"
	"grpc/stream/calculator/calculatorpb"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage  started \n")
	sum := int32(0)
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: float64(sum) / float64(count),
			})

		}
		if err != nil {
			log.Fatalf("Error with compute average recv : %v\n", err)
		}
		sum += req.GetNumber()
		count++
	}

}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Printf("SquareRoot  started \n")
	num := req.GetNumber()
	if num < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Received an invalid negative number : %v\n", num))
	}
	return &calculatorpb.SquareRootResponse{NumberRoot: math.Sqrt(float64(num))}, nil
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("FindMaximum  started \n")
	curMax := 0
	for {
		req, err := stream.Recv()
		if curMax < int(req.GetNumber()) {
			curMax = int(req.GetNumber())
			err = stream.Send(&calculatorpb.FindMaximumResponse{
				Number: int32(curMax),
			})
			if err != nil {
				log.Printf("Got an error during sendinging : %v", err)

			}
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Error with Find maximum recv : %v\n", err)
			return err
		}
	}
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
