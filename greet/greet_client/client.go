package main

import (
	"context"
	"fmt"
	"gorpc-go-course/greet/greetpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("hello I am a client")
	cc, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect : %v", err)
	}
	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)
	doUnary(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Printf("starting unary rpc client: %v\n", c)
	req := greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ramesh",
			LastName:  "Naidu",
		},
	}
	res, err := c.Greet(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet : %v\n", res.Result)

}
