package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"warehouse-cli/database/database"
	di "warehouse-cli/database/databaseInstance"
	"warehouse-cli/handlers"
)

func main() {

	ReplicationFactor := 0

	http.HandleFunc("/GET", handlers.GET)
	http.HandleFunc("/SET", handlers.SET)
	http.HandleFunc("/DELETE", handlers.DELETE)
	http.HandleFunc("/home", handlers.Welcome)

	database.StartListeningForChanges()

	di.Instances = make([]di.DatabaseInstance, 0, 3)
	di.CreateNewDatabaseInstance("127.0.0.1", 8765, &ReplicationFactor)
	di.CreateNewDatabaseInstance("127.0.0.1", 9876, &ReplicationFactor)
	di.CreateNewDatabaseInstance("127.0.0.1", 8697, &ReplicationFactor)
	di.CreateNewDatabaseInstance("127.0.0.1", 8080, &ReplicationFactor)

	// Prevent the main goroutine from exiting
	for {
		eventLoop(&ReplicationFactor)
		select {}
	}

}

func eventLoop(replicationFactor *int) {
	go func() {
		for {
			var input string
			_, err := fmt.Scanln(&input)
			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}
			switch input {
			case "start":
				host, port, err := parseInput()
				if err != nil {
					continue
				}
				log.Print(host)
				log.Print(port)
				di.CreateNewDatabaseInstance(host, port, replicationFactor)
			case "kill":
				host, port, err := parseInput()
				if err != nil {
					continue
				}
				di.DeleteDatabaseInstance(host, port, replicationFactor)
			default:
				fmt.Println("Invalid command:", input)
			}
		}
	}()
}

func parseInput() (host string, port int, err error) {

	fmt.Println("enter host")
	_, err = fmt.Scan(&host)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return "", 0, err
	}
	fmt.Println("enter port")
	var portS string
	_, err = fmt.Scanln(&portS)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return "", 0, err
	}
	port, err = strconv.Atoi(portS)
	if err != nil {
		fmt.Println("invalid port")
		return "", 0, err
	}
	return host, port, nil
}
