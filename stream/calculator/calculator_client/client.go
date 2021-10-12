package main

import (
	"context"
	"errors"
	"fmt"
	"grpc/stream/calculator/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// doClientStreaming(c)
	// doBidiStreaming(c)
	doErrorUnary(c)

}
func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a SQRT error RPC .. ")
	num := int32(-9)
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: num})
	if err != nil {
		err, ok := status.FromError(err)
		if ok {
			fmt.Println(err.Message(), err.Code())
			if err.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a -ve number")
				return
			}
		} else {
			log.Fatalf("error calling square root : %v", err)
		}
	}
	fmt.Printf("Result of square root of %v: %v", num, res.GetNumberRoot())

}

func doBidiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Bidi streaming RPC .. ")
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream and calling FindMaximum : %v", err)
	}
	wc := make(chan int)
	go func() {
		numbers := []int32{4, 5, 53, 6, 7, 756}
		for _, number := range numbers {
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(time.Second)
		}
		stream.CloseSend()
	}()
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				wc <- 0
				break
			}
			if err != nil {
				log.Printf("Error while receinving : %v", err)
				wc <- 1
			}
			maximum := res.GetNumber()
			log.Printf("Got the new maximum : %v", maximum)
		}
	}()
	<-wc
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

func BSTSearch(root *TreeNode, key int) (*TreeNode, *TreeNode, *TreeNode) {
	if root == nil {
		return nil, nil, nil
	}
	curr := root
	var succ, pred *TreeNode
	for curr != nil {
		if key == curr.key {
			return curr, succ, pred
		} else if key < curr.key {
			succ = curr
			curr = curr.left
		} else {
			pred = curr
			curr = curr.right
		}
	}
	return nil, succ, pred
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
		return nil, nil //empty tree
	}
	curr := root
	for curr.left != nil {
		curr = curr.left
	}
	return curr, nil
}

func BSTMax(root *TreeNode) (*TreeNode, error) {
	if root == nil {
		return nil, nil //empty tree
	}
	curr := root
	for curr.right != nil {
		curr = curr.right
	}
	return curr, nil
}

func BSTSuccessor(root *TreeNode, key int) (*TreeNode, error) {
	if root == nil {
		return nil, nil //empty tree
	}
	curr, succ, _ := BSTSearch(root, key)
	if curr == nil {
		return nil, errors.New("key not found")
	}
	if curr.right != nil {
		//leftmost child of the right subtree
		return BSTMin(curr.right)
	}
	//now we have to go back to the ancestral tree to find the first right turn
	return succ, nil
}

func BSTPredecessor(root *TreeNode, key int) (*TreeNode, error) {
	if root == nil {
		return nil, nil //empty tree
	}
	curr, _, pred := BSTSearch(root, key)
	if curr == nil {
		return nil, errors.New("key not found")
	}
	if curr.left != nil {
		//leftmost child of the right subtree
		return BSTMax(curr.left)
	}
	//now we have to go back to the ancestral tree to find the first right turn
	return pred, nil
}
