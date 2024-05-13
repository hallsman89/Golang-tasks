package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	arr := getRandomData(10)
	fmt.Println(arr)

	ch := sleepSort(arr)
	for i := 0; i < len(arr); i++ {
		select {
		case val, ok := <-ch:
			if ok {
				fmt.Println(val)
			} else {
				fmt.Println("Channel closed unexpectedly")
				return
			}
		case <-time.After(10 * time.Second):
			fmt.Println("Timeout occurred")
			return
		}
	}
}

func sleepSort(arr []int) chan int {
	ch := make(chan int, len(arr))
	var wg sync.WaitGroup
	for i := 0; i < len(arr); i++ {
		val := arr[i]
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			time.Sleep(time.Duration(val) * time.Second)
			ch <- val
		}(val)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch
}

func getRandomData(cap int) []int {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	res := make([]int, cap)
	for i := 0; i < cap; i++ {
		res[i] = r.Intn(9)
	}
	return res
}
