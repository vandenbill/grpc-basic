package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/vandenbill/grpc-basic/calculator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	cc, err := grpc.Dial("0.0.0.0:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not dial: %v\n", err)
	}

	c := pb.NewCalculatorServiceClient(cc)
	// biDirectional(c)
	unary(c)
}

func unary(c pb.CalculatorServiceClient) {
	r, err := c.Calculate(context.Background(), &pb.CalculatorRequest{V1: -1, V2: 323232})
	if err != nil {
		log.Fatalf("Could not calculate: %v", err)
	}
	if err != nil {
		status, ok := status.FromError(err)
		if ok {
			log.Printf("Code Message: %s\n", status.Message())
			log.Printf("Code error: %s\n", status.Message())

			if status.Code() == codes.InvalidArgument {
				log.Printf("Could not send negative\n")
			}
		}
	}
	log.Printf("Result: %d", r.R)
}

func serverStream(c pb.CalculatorServiceClient) {
	stream, err := c.CalculateManyTimes(context.Background(), &pb.CalculatorRequest{V1: 1, V2: 2})
	if err != nil {
		log.Fatalf("Error while calling CalculateManyTimes %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while read stream %v", err)
		}
		log.Printf("Result: %d", msg.GetR())
	}
}

func clientStream(c pb.CalculatorServiceClient) {
	requests := []*pb.CalculatorRequest{
		&pb.CalculatorRequest{
			V1: 1,
			V2: 2,
		}, &pb.CalculatorRequest{
			V1: 3,
			V2: 4,
		},
	}

	stream, err := c.CalculateLongClientStream(context.Background())
	if err != nil {
		log.Fatalf("Error while calling CalculateLongClientStream: %v\n", err)
	}

	for i, v := range requests {
		log.Printf("Sending request: %d", i)
		stream.Send(v)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v\n", err)
	}

	log.Print("Response: ", res.R)
}

func biDirectional(c pb.CalculatorServiceClient) {
	requests := []*pb.CalculatorRequest{
		&pb.CalculatorRequest{
			V1: 1,
			V2: 2,
		}, &pb.CalculatorRequest{
			V1: 3,
			V2: 4,
		},
	}

	stream, err := c.CalculateBiDirectional(context.Background())
	if err != nil {
		log.Fatalf("Error while calling CalculateBiDirectional: %v\n", err)
	}

	chanl := make(chan struct{})

	go func() {
		for i, v := range requests {
			log.Printf("Sending request: %d", i)
			stream.Send(v)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receive response: %v\n", err)
			}
			log.Printf("Result: %d", req.GetR())
		}
		close(chanl)
	}()

	<-chanl
}
