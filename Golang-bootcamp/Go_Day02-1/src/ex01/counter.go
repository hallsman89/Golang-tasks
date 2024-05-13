package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func countThings(fileName string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(fileName)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	var l, m, w int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l++
		m += len(scanner.Text())
		w += len(strings.Fields(scanner.Text()))
	}
	m += l

	printThing(l, m, w, fileName)
}

func printThing(l, m, w int, f string) {
	printLock.Lock()
	defer printLock.Unlock()

	if fl.l {
		fmt.Printf("%d\t%s\n", l, f)
	} else if fl.m {
		fmt.Printf("%d\t%s\n", m, f)
	} else if fl.w {
		fmt.Printf("%d\t%s\n", w, f)
	}
}

var printLock sync.Mutex
