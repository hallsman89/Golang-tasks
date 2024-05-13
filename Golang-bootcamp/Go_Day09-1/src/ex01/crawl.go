package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
)

const MaxConcurrency = 8

func crawlWeb(ctx context.Context, urls <-chan string) chan *string {
	results := make(chan *string)
	go runRoutines(ctx, urls, results)
	return results
}

func runRoutines(ctx context.Context, urls <-chan string, res chan<- *string) {
	semaphore := make(chan struct{}, MaxConcurrency)
	wg := &sync.WaitGroup{}

	for url := range urls {
		select {
		case <-ctx.Done():
			return
		case semaphore <- struct{}{}:
			wg.Add(1)
			go func(u string) {
				defer wg.Done()
				body, err := fetchURL(u)
				if err == nil {
					res <- &body
				}
				<-semaphore
			}(url)
		}
	}

	wg.Wait()
	close(res)
}

func fetchURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-OK status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
