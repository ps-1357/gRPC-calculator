package main

import (
	"calculator/calculatorpb"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Calculator client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorClient(cc)
	//fmt.Printf("Created client: %f", c)
	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBidiStreaming(c)

	doErrorUnary(c)
}
func doUnary(c calculatorpb.CalculatorClient) {
	fmt.Println("Starting to do a Sum Unary RPC")
	req := &calculatorpb.SumRequest{
		FirstNumber:  5,
		SecondNumber: 40,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", res.SumResult)
}

func doServerStreaming(c calculatorpb.CalculatorClient) {
	fmt.Println("Starting to do a Prime Decomposition Server Streaming RPC")
	req := &calculatorpb.PrimeNumberDecompostionRequest{
		Number: 145890,
	}
	stream, err := c.PrimeNumberDecompostion(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}
}

func doClientStreaming(c calculatorpb.CalculatorClient) {
	fmt.Println("Starting to do a ComputeAverage Client Streaming RPC...")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}
	numbers := []int64{3, 5, 9, 54, 60, 71}
	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}
	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while receinving response: %v", err)
	}
	fmt.Printf("The Average is: %v", res.GetAverage())
}

func doBidiStreaming(c calculatorpb.CalculatorClient) {
	fmt.Println("Starting to do a FindMaximum Bidi Streaming RPC ...")

	stream, err := c.FindMaximum(context.Background())

	if err != nil {
		log.Fatalf("Error while opening stream and calling FindMaximum: %v\n", err)
	}

	waitc := make(chan struct{})

	//send go routine
	go func() {
		numbers := []int32{4, 7, 2, 19, 4, 6, 32}
		for _, number := range numbers {
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep((1000 * time.Millisecond))
		}
		stream.CloseSend()
	}()
	//receive go routine
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Problem while reading server stream: %v\n", err)
				break
			}
			maximum := res.GetMaximum()
			fmt.Printf("Received a new maximum of ...: %v\n", maximum)
		}
		close(waitc)
	}()
	<-waitc
}

func doErrorUnary(c calculatorpb.CalculatorClient) {
	fmt.Println("Starting to do a SquareRoot Unary RPC...")
	// correct call
	doErrorCall(c, 10)

	// error call
	doErrorCall(c ,-2)
	
}

func doErrorCall( c calculatorpb.CalculatorClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			fmt.Printf("Error message from server: %v\n",respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number")
				return
			}
		} else {
			log.Fatalf("Big Error calling SquareRoot: %v", err)
			return
		}
	}
	fmt.Printf("Result of sqare root of %v: %v\n", n, res.GetNumberRoot())
	// error call
}	
