package main

import (
	pb "ex02/data/go"
	"ex02/repository"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"math"
	"math/rand"
	"net"
	"time"
)

type Server struct {
	pb.UnimplementedMilitaryDeviceServer
	repository *repository.Repository
}

func (s *Server) Connect(params *pb.ConnectionParams, stream pb.MilitaryDevice_ConnectServer) error {
	num := 0
	for {

		message := &pb.Message{
			SessionId: uuid.New().String(),
			Frequency: generateFrequency(params.Mean, params.StdDeviation),
			Timestamp: timestamppb.Now(),
		}
		s.repository.POST(message)

		if err := stream.Send(message); err != nil {
			log.Print(err)
			return err
		}
		log.Printf("message number %d sent\n", num)
		num++

		time.Sleep(time.Second)
	}
}

func generateFrequency(mean, stdDeviation float64) float64 {
	u1 := rand.Float64()
	u2 := rand.Float64()
	z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)

	frequency := z0*stdDeviation + mean

	return frequency
}

var (
	port = flag.Int("port", 8080, "port to listen on")
)

func init() {
	flag.Parse()
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer func() { _ = lis.Close() }()

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterMilitaryDeviceServer(grpcServer, &Server{})
	repo := &repository.Repository{}
	repo.Connect()

	grpcServer = grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterMilitaryDeviceServer(grpcServer, &Server{repository: repo})

	log.Printf("Server started listening on port %d", *port)
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
