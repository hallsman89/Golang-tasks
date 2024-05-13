package main

import (
	"day02/ex03/archiver"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

var pathToDir *string

func init() {
	pathToDir = flag.String("a", ".", "path to directory")
	flag.Parse()
	if err := isCorrectFolder(*pathToDir); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	files, err := getFiles()
	if err != nil {
		log.Fatalln(err)
	}

	wg := new(sync.WaitGroup)
	for _, file := range files {
		wg.Add(1)
		go archiver.ArchiveFile(*pathToDir, file, wg)
	}
	wg.Wait()
}

func isCorrectFolder(dir string) error {
	if dir == "" {
		return errors.New("error: empty directory path")
	}
	info, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("path %s must point to a directory", dir)
	}
	return nil
}

func getFiles() ([]string, error) {
	var files []string
	args := flag.Args()
	if len(args) == 0 {
		return nil, errors.New("error: need to specify 1 or more files")
	}
	if *pathToDir == "." {
		files = args
	} else {
		files = args[2:]
	}

	return files, nil
}
