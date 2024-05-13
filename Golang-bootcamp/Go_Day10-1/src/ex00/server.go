package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	pb "data/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type MiliServiceServer struct {
	pb.UnimplementedMiliServiceServer
}

func (s *MiliServiceServer) SendMessage(ctx context.Context, in *pb.Message) (*pb.InfoAboutConnection, error) {
	sessionID := in.GetSessionId()
	log.Printf("New connection with session ID: %s\n", sessionID)
	mean := rand.Float64()*20 - 10     // [-10, 10]
	stdDev := rand.Float64()*1.2 + 0.3 // [0.3, 1.5]
	log.Printf("Mean: %f, StdDev: %f\n", mean, stdDev)
	return &pb.InfoAboutConnection{
		SessionId:    sessionID,
		Mean:         mean,
		StdDeviation: stdDev,
	}, nil
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
