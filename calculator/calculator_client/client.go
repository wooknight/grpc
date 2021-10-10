package main

import (
	"context"
	"fmt"
	"gorpc-go-course/calculator/calculatorpb"
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
	doUnary(c)
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
