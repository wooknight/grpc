package main

import (
	"context"
	"errors"
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
	// doServerStreaming(c)
	doClientStreaming(c)

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

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a client streaming for average")
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream : %v\n", err)
	}
	numbers := []int32{3, 5, 9, 54, 23}
	for _, num := range numbers {
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: num,
		})
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while reciving response  : %v\n", err)
	}
	fmt.Printf("The average is %v", res.GetAverage())
}

type TreeNode struct {
	key         int
	value       int
	left, right *TreeNode
}

func BSTSearch(root *TreeNode, key int) *TreeNode {
	if root == nil {
		return nil
	}
	curr := root
	for curr != nil {
		if key == curr.key {
			return curr
		} else if key < curr.key {
			curr = curr.left
		} else {
			curr = curr.right
		}
	}
	return nil
}

func BSTInsert(root *TreeNode, key int, val int) (*TreeNode, error) {
	tmp := TreeNode{key: key, value: val}
	if root == nil {
		return &tmp, nil //empty tree so we are creating the new node
	}
	curr := root

	var prev *TreeNode
	for curr != nil {
		if key == curr.key {
			return nil, errors.New("Key already exists")
		} else if key < curr.key {
			prev = curr
			curr = curr.left
		} else {
			prev = curr
			curr = curr.right
		}
	}
	if key < prev.key {
		prev.left = &tmp
	} else {
		prev.right = &tmp
	}
	return root, nil
}

func BSTMin(root *TreeNode) (*TreeNode, error) {
	if root == nil {
		return nil, nil //empty tree so we are creating the new node
	}
	curr := root
	for curr.left != nil {
		curr = curr.left
	}
	return curr, nil
}

func BSTMax(root *TreeNode) (*TreeNode, error) {
	if root == nil {
		return nil, nil //empty tree so we are creating the new node
	}
	curr := root
	for curr.right != nil {
		curr = curr.right
	}
	return curr, nil
}
