package main

import (
	"context"
	"log"
	"math"
	"net"
	"os"
	"sync"

	pb "data/gen/go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type MiliServiceServer struct {
	pb.UnimplementedMiliServiceServer

	mean      float64
	stdDev    float64
	valueList []float64
	valueMu   sync.Mutex
}

func (s *MiliServiceServer) SendMessage(ctx context.Context, in *pb.Message) (*pb.InfoAboutConnection, error) {
	sessionID := in.GetSessionId()
	log.Printf("New connection with session ID: %s\n", sessionID)
	s.processValues(in.GetFrequency())
	s.valueMu.Lock()
	mean := s.mean
	stdDev := s.stdDev
	s.valueMu.Unlock()
	log.Printf("Mean: %f, StdDev: %f\n", mean, stdDev)

	return &pb.InfoAboutConnection{
		SessionId:    sessionID,
		Mean:         mean,
		StdDeviation: stdDev,
	}, nil
}

func (s *MiliServiceServer) processValues(value float64) {
	s.valueMu.Lock()
	defer s.valueMu.Unlock()
	s.valueList = append(s.valueList, value)
	if len(s.valueList) > 50 {
		s.valueList = s.valueList[1:]
	}
	if len(s.valueList) > 50 && s.stdDev < 0.1 {
		log.Printf("Predicted parameters are considered accurate. Mean: %f, StdDev: %f\n",
			s.mean, s.stdDev)
	}
	if len(s.valueList) > 50 && s.stdDev > 0.1 {
		log.Printf("Anomaly detected. Stopping the server.\n")
		os.Exit(1)
	}
	sum := 0.0
	for _, v := range s.valueList {
		sum += v
	}
	s.mean = sum / float64(len(s.valueList))
	sumSquares := 0.0
	for _, v := range s.valueList {
		sumSquares += (v - s.mean) * (v - s.mean)
	}
	s.stdDev = math.Sqrt(sumSquares / float64(len(s.valueList)))
	if len(s.valueList)%10 == 0 {
		log.Printf("Processed %d values. Mean: %f, StdDev: %f\n", len(s.valueList), s.mean, s.stdDev)
	}
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMiliServiceServer(s, &MiliServiceServer{})
	reflection.Register(s)
	log.Println("Server is starting...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
