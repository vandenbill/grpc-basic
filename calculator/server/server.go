package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	pb "github.com/vandenbill/grpc-basic/calculator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var Addr string = "0.0.0.0:8081"

type Server struct {
	pb.CalculatorServiceServer
}

func main() {
	lis, err := net.Listen("tcp", Addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Listen to: %s", Addr)

	s := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(s, &Server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *Server) Calculate(ctx context.Context, in *pb.CalculatorRequest) (*pb.CalculatorResponse, error) {
	log.Printf("Calculate func was invoked")
	if in.GetV1() < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprint("Could'n be negative"))
	}
	return &pb.CalculatorResponse{
		R: in.V1 + in.V2,
	}, nil
}

func (s *Server) CalculateManyTimes(in *pb.CalculatorRequest, stream pb.CalculatorService_CalculateManyTimesServer) error {
	log.Printf("CalculateManyTimes func was invoked")

	for i := 0; i < 10; i++ {
		stream.Send(&pb.CalculatorResponse{R: int64(int(in.V1) * int(in.V2) * i)})
	}

	return nil
}

func (s *Server) CalculateLongClientStream(stream pb.CalculatorService_CalculateLongClientStreamServer) error {
	log.Printf("CalculateLongClientStream func was invoked")

	result := 0
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&pb.CalculatorResponse{R: int64(result)})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		result += int(req.V1)
		result += int(req.V2)
	}
}

func (s *Server) CalculateBiDirectional(stream pb.CalculatorService_CalculateBiDirectionalServer) error {
	log.Printf("CalculateBiDirectional func was invoked")

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		err = stream.Send(&pb.CalculatorResponse{R: int64(math.Pow(float64(req.V1), float64(req.V2)))})
		if err != nil {
			log.Fatalf("Error while sending data to client: %v", err)
			return err
		}
	}
}
