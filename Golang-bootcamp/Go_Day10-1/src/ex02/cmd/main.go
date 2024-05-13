package main

import (
	pb "ex02/data/go"
	"ex02/repository"
	"flag"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math"
	"math/rand"
)

var (
	anomaly = flag.Float64("k", 1.5, "anomaly coefficient")
)

const (
	MIN_DEVIATION = -10
	MAX_DEVIATION = 10
	MIN_SD        = 0.3
	MAX_SD        = 1.5
	MAX_LOG_SIZE  = 50
	ADDR          = "localhost:8080"
)

func getMean(frequencies []float64) float64 {
	sum := 0.0
	for _, freq := range frequencies {
		sum += freq
	}
	return sum / float64(len(frequencies))
}

func calculateStandardDeviation(frequencies []float64, mean float64) float64 {
	sumSquaredDiff := 0.0

	for _, freq := range frequencies {
		diff := freq - mean
		sumSquaredDiff += diff * diff
	}

	variance := sumSquaredDiff / float64(len(frequencies))

	return math.Sqrt(variance)
}

func main() {
	flag.Parse()

	r := repository.Repository{}
	r.Connect()

	conn, err := grpc.Dial(ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to %s: %v", ADDR, err)
	}
	defer func() { _ = conn.Close() }()
	client := pb.NewMilitaryDeviceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	request, err := client.Connect(ctx, &pb.ConnectionParams{
		Mean:         rand.Float64()*(MAX_DEVIATION*2) + MIN_DEVIATION,
		StdDeviation: rand.Float64()*(MAX_SD-MIN_SD) + MIN_SD,
	})
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	freqs := make([]float64, 0, MAX_LOG_SIZE)
	i := 1
	for ; i <= MAX_LOG_SIZE; i++ {
		msg, err := request.Recv()
		if err != nil {
			log.Fatalf("failed to receive msg: %v", err)
		}
		r.POST(msg)
		if i%10 == 0 {
			mean := getMean(freqs)
			sd := calculateStandardDeviation(freqs, mean)
			log.Printf("Processed %d messages\nMean is: %f\nSTD is: %f", i, mean, sd)
		}
		freqs = append(freqs, msg.Frequency)
	}
	mean := getMean(freqs)
	sd := calculateStandardDeviation(freqs, mean)
	acceptedRange := *anomaly * sd
	for {
		msg, err := request.Recv()
		if err != nil {
			log.Fatalf("failed to receive msg: %v", err)
		}
		r.POST(msg)
		if msg.Frequency > mean+acceptedRange || msg.Frequency < mean-acceptedRange {
			log.Printf("message %d is an anomaly.\nFrequency: %f, expected: +-%f ",
				i, msg.Frequency, acceptedRange)
				r.POST(msg)
		}
		i++
	}
}
