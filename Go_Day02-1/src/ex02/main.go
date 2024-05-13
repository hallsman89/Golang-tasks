package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("Usage: myXargs <command>")
	}

	console := bufio.NewScanner(os.Stdin)
	for console.Scan() {
		input := console.Text()
		if input == "" {
			continue
		}

		cmd := exec.Command(args[0], append(args[1:], input)...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Error executing command: %s\n", err)
		}
	}

	if err := console.Err(); err != nil {
		log.Fatal(err)
	}
}
