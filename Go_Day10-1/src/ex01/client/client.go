package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	pb "data/gen/go"
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
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	anomalyCoefficient := flag.Float64("k", 2.0, "STD anomaly coefficient")
	flag.Parse()

	for i := 0; i < 100; i++ {
		message := &pb.Message{
			SessionId:    "example_session_id",
			Frequency:    rand.Float64() * 10,
			TimestampUtc: time.Now().Unix(),
		}
		response, err := client.SendMessage(ctx, message)
		if err != nil {
			log.Fatalf("Error calling SendMessage: %v", err)
		}
		fmt.Printf("Response from the server: %+v\n", response)
		if message.Frequency > response.Mean+(response.StdDeviation**anomalyCoefficient) ||
			message.Frequency < response.Mean-(response.StdDeviation**anomalyCoefficient) {
			log.Printf("Anomaly detected! Frequency: %f. Stopping the client.\n", message.Frequency)
			os.Exit(1)
		}
		time.Sleep(time.Second)
	}
}
