package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Flags struct {
	f, d, sl bool
	ext      string
}

var fl Flags

func init() {
	flag.BoolVar(&fl.f, "f", false, "file")
	flag.BoolVar(&fl.d, "d", false, "directory")
	flag.BoolVar(&fl.sl, "sl", false, "symbolic links")
	flag.StringVar(&fl.ext, "ext", "", "extension")
}

func main() {
	flag.Parse()

	filePath, err := checkFlags()
	if err != nil {
		log.Fatal(err)
	}

	err = filepath.Walk(filePath, WalkFunc)
	if err != nil {
		log.Fatal(err)
	}
}

func checkFlags() (string, error) {
	args := flag.Args()
	if len(args) == 0 {
		return "", errors.New("ERROR: missing file path")
	}

	if !fl.f && fl.ext != "" {
		fmt.Println("WARNING: -ext may only be used with -f")
	}

	if !fl.f && !fl.d && !fl.sl {
		fl.f, fl.d, fl.sl = true, true, true
	}

	// Check if the file path exists
	filePath := args[0]
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", errors.New("ERROR: file path does not exist")
	}

	return filePath, nil
}
