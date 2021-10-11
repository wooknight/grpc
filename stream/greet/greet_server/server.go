package main

import (
	"context"
	"fmt"
	"grpc/stream/calculator/calculatorpb"
	"grpc/stream/greet/greetpb"
	"io"
	"time"

	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked . req : %v", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	result := "Hello " + firstName + "  " + lastName
	res := greetpb.GreetResponse{
		Result: result,
	}
	return &res, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	firstName := req.GetGreeting().GetFirstName()
	result := "First hello " + firstName + " " + time.Now().String()
	log.Println(result)
	for i := 0; i < 10; i++ {
		<-time.After(10 * time.Second)
		result = "hello " + firstName + "  " + time.Now().String()

		res := greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(&res)
	}
	return nil
}

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	result := "Long Greet was invoked " + time.Now().String()
	log.Println(result)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading long greet stream : %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result += firstName + "! \n"
	}
	return stream.SendAndClose(&greetpb.LongGreetResponse{
		Result: result,
	})
}

func main() {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}

}
