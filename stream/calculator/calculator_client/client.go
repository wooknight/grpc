package main

import (
	"context"
	"fmt"
	"grpc/stream/calculator/calculatorpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I am a calc client ")
	cc, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()
	c := calculatorpb.NewCalculatorServiceClient(cc)
	// doUnary(c)
	doServerStreaming(c)

}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a SUM unary RPC .. ")
	req := calculatorpb.SumRequest{
		FirstNumber:  25,
		SecondNumber: 156,
	}
	res, err := c.Sum(context.Background(), &req)
	if err != nil {
		log.Fatalf("error while calling SUM RPC %v", err)
	}
	log.Printf("Response from Sum : %v", res.SumResult)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a server streaming for prime number decomposition")
	req := calculatorpb.PrimeNumberDecompositionRequest{
		Number: 1233673,
	}
	resStream, err := c.PrimeNumberDecomposition(context.Background(), &req)
	if err != nil {
		log.Fatalf("error while trying to get stream from prime decomposition: %v", err)
	}
	for {
		res, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while calling decomposition rpc: %v", err)
		}
		log.Printf("Response from Prime number decomposition : %v", res.GetPrimeFactor())
	}

}
