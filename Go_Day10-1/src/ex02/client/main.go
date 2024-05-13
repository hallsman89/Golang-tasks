package main

import (
	"context"
	"io"
	"log"

	pb "ex02/data/go"
	"ex02/repository"
	"google.golang.org/grpc"
)

const (
	serverAddress = "localhost:8080"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer func() { _ = conn.Close() }()

	client := pb.NewMilitaryDeviceClient(conn)

	// Инициализация репозитория
	repo := &repository.Repository{}
	repo.Connect()

	// Отправка запроса на сервер
	params := &pb.ConnectionParams{
		Mean:         10.0,
		StdDeviation: 2.0,
	}

	stream, err := client.Connect(context.Background(), params)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error receiving message: %v", err)
		}

		// Сохранение в репозиторий
		repo.POST(message)

		log.Printf("Received message: %+v", message)
	}
}
