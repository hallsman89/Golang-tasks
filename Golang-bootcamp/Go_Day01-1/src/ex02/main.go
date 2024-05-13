package main

import (
	"bufio"
	"fmt"
	"flag"
	"log"
	"os"
)

type Flags struct {
	OldFile, NewFile string
}

var fl Flags

func init() {
	flag.StringVar(&fl.OldFile, "old", "", "original db")
	flag.StringVar(&fl.NewFile, "new", "", "original db")
}

func main() {
	flag.Parse()
	if fl.OldFile == "" || fl.NewFile == "" {
		fmt.Println("ERROR: Please provide both --old and --new file names")
		return
	}

	compareDumps(fl.OldFile, fl.NewFile)
}

func compareDumps(fileName1, fileName2 string) {
    data1 := getFileMap(fileName1)
    data2 := getFileMap(fileName2)

    for key, _ := range data2 {
        if _, ok := data1[key]; !ok {
            fmt.Println("ADDED:", key)
        }
    }

    for key, _ := range data1 {
        if _, ok := data2[key]; !ok {
            fmt.Println("REMOVED:", key)
        }
    }
}


func getFileMap(fileName string) map[string]bool {
    file, err := os.Open(fileName)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    lines := make(map[string]bool)
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        lines[line] = true
    }
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    return lines
}


