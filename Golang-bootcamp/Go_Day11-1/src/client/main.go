package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Address struct {
	host string
	port string
}

func (a *Address) String() string {
	return fmt.Sprintf("%s:%s", a.host, a.port)
}

var address Address

func init() {
	flag.StringVar(&address.host, "H", "", "host")
	flag.StringVar(&address.port, "P", "", "port")
}

func main() {
	flag.Parse()

	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/home", address.String()), nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(responseBody))

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := scanner.Text()
		line := strings.Split(text, " ")
		if len(line) >= 2 {
			err := reqProces(client, line[0], line[1], strings.Join(line[2:], " "))
			if err != nil {
				log.Fatal(err)
			}
		}

	}
}
func reqProces(client *http.Client, method, key, body string) error {
	var reqURL string

	if method == "SET" {
		reqURL = fmt.Sprintf("http://%s/%s?uuid=%s&body=%s", address.String(), method, key, url.QueryEscape(body))
	} else {
		reqURL = fmt.Sprintf("http://%s/%s?uuid=%s", address.String(), method, key)
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(responseBody))

	return nil
}
