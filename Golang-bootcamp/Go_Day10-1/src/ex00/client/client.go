package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "data/gen/go" // Импортируем сгенерированный пакет

	"google.golang.org/grpc"
)

func main() {
	address := "localhost:8080"
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to establish connection: %v", err)
	}
	defer conn.Close()
	client := pb.NewMiliServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	message := &pb.Message{
		SessionId:    "example_session_id",
		Frequency:    3.14,
		TimestampUtc: 123456789,
	}
	response, err := client.SendMessage(ctx, message)
	if err != nil {
		log.Fatalf("Error calling SendMessage: %v", err)
	}
	fmt.Printf("Response from the server: %+v\n", response)
}
