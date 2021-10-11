package main

import (
	"context"
	"fmt"
	"grpc/stream/greet/greetpb"
	"io"
	"log"
	"time"

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
	// doServerStreaming(c)
	// doClientStreaming(c)
	doBidiStreaming(c)
}

func doBidiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Doing a client BIDI streaming request")
	requests := []*greetpb.LongGreetManyRequest{
		&greetpb.LongGreetManyRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ramesh",
			},
		},
		&greetpb.LongGreetManyRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Vanguard",
			},
		},
		&greetpb.LongGreetManyRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Hannibal",
			},
		},
		&greetpb.LongGreetManyRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Seneca",
			},
		},
		&greetpb.LongGreetManyRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bhosada",
			},
		},
	}
	stream, err := c.LongGreetMany(context.Background())
	if err != nil {
		log.Printf("Error while creating stream: %v", err)
		return
	}
	wc := make(chan int)
	go func() {
		for _, req := range requests {
			log.Printf("Sending request : %v\n", req)
			stream.Send(req)
			time.Sleep(time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				wc <- 0
				return
			}
			if err != nil {
				log.Printf("Error while receiving : %v", err)
				close(wc)
				return
			}
			fmt.Printf("Received : %v\n", res.GetResult())
		}
	}()

	<-wc
}
func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Doing a client streaming request")
	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ramesh",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Vanguard",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Hannibal",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Seneca",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bhosada",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while getting long Greet : %v", err)
	}
	for _, req := range requests {
		fmt.Printf("Sending req : %s\n", req.GetGreeting().GetFirstName())
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Received an error while closing : %v\n", err)
	}
	fmt.Printf("longGreetResponse ; %v\n", res)
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
