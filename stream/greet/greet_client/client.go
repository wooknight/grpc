package main

import (
	"context"
	"fmt"
	"grpc/stream/greet/greetpb"
	"io"
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
	// doUnary(c)
	doServerStreaming(c)
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

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting server streaming RPC....")

	req := greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ramesh",
			LastName:  "Naidu",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), &req)
	if err != nil {
		log.Fatalf("Client failed to initiate greet many times request. Error : %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while receiving. Error : %v", err)
		}
		log.Printf("Response from GreetManyTimes : %v", msg.GetResult())
	}

}
