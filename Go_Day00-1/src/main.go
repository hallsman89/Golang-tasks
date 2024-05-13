package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/montanaflynn/stats"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	parser := bufio.NewScanner(os.Stdin)
	for parser.Scan() {
		var slice []int
		strSlice := strings.Fields(parser.Text())
		for _, elem := range strSlice {
			if value, err := strconv.Atoi(elem); err == nil {
				slice = append(slice, value)
			} else {
				log.Fatal("Bad data")
			}
		}

		if len(os.Args) == 1 {
			printAllStats(slice)
		} else {
			flagSet := flag.NewFlagSet("", flag.ExitOnError)
			meanFlag := flagSet.Bool("mean", false, "Calculate mean")
			medianFlag := flagSet.Bool("median", false, "Calculate median")
			modeFlag := flagSet.Bool("mode", false, "Calculate mode")
			sdFlag := flagSet.Bool("sd", false, "Calculate standard deviation")

			flagSet.Parse(os.Args[1:])

			if *meanFlag {
				fmt.Print("Mean: ")
				mean(slice)
			}
			if *medianFlag {
				fmt.Print("Median: ")
				median(slice)
			}
			if *modeFlag {
				fmt.Print("Mode: ")
				mode(slice)
			}
			if *sdFlag {
				fmt.Print("Standard Deviation: ")
				sd(slice)
			}
			fmt.Println()
		}
	}
}

func printAllStats(slice []int) {
	fmt.Println("Mean:")
	mean(slice)
	fmt.Println("Median:")
	median(slice)
	fmt.Println("Mode:")
	mode(slice)
	fmt.Println("Standard Deviation:")
	sd(slice)
}

func mean(slice []int) {
	data := stats.LoadRawData(slice)
	result, err := stats.Mean(data)
	if err != nil {
		fmt.Println("Bad data")
	} else {
		fmt.Printf("%.2f\n", result)
	}
}

func median(slice []int) {
	data := stats.LoadRawData(slice)
	result, err := stats.Median(data)
	if err != nil {
		fmt.Println("Bad data")
	} else {
		fmt.Printf("%.2f\n", result)
	}
}

func mode(slice []int) {
	data := stats.LoadRawData(slice)
	result, err := stats.Mode(data)
	if err != nil {
		fmt.Println("Bad data")
	} else {
		if len(result) == 0 {
			min, _ := stats.Min(data)
			fmt.Printf("%.2f\n", min)
		} else if len(result) > 1 {
			min, _ := stats.Min(result)
			fmt.Printf("%.2f\n", min)
		} else {
			fmt.Printf("%.2f\n", result[0])
		}
	}
}

func sd(slice []int) {
	data := stats.LoadRawData(slice)
	result, err := stats.StandardDeviation(data)
	if err != nil {
		fmt.Println("Bad data")
	} else {
		fmt.Printf("%.2f\n", result)
	}
}
