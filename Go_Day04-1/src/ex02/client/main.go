package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	candyType := flag.String("k", "AA", "two-letter abbreviation for the candy type")
	candyCount := flag.Int("c", 1, "count of candy to buy")
	money := flag.Int("m", 20, "amount of money you gave to machine")
	flag.Parse()

	client := getClient()

	bodyResponse := struct {
		Money      int `json:"money"`
		CandyType  string `json:"candyType"`
		CandyCount int `json:"candyCount"`
	}{
		Money:      *money,
		CandyType:  *candyType,
		CandyCount: *candyCount,
	}

	bodyJSON, _ := json.Marshal(bodyResponse)

	resp, err := client.Post("https://localhost:8080/buy_candy", "application/json", bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %s Body: %s\n", resp.Status, string(body))
}

func getClient() *http.Client {
	data, _ := os.ReadFile("ca/minica.pem")
	cp, _ := x509.SystemCertPool()
	cp.AppendCertsFromPEM(data)

	config := &tls.Config{
		RootCAs:               cp,
		GetClientCertificate:  nil,
		VerifyPeerCertificate: nil,
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: config,
		},
	}

	return client
}