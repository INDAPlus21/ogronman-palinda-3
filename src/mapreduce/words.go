package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

const DataFile = "loremipsum.txt"


// Return the word frequencies of the text argument.
func WordCount(text string) map[string]int {

	freqs := make(map[string]int)
	words := strings.Fields(text)
	readers := 5
	amount := len(words)
	worksize :=  amount / readers
	var wg sync.WaitGroup

	freqsCh := make(chan map[string]int, worksize/readers+1)

	for i, j := 0, worksize; i < amount; i, j = j, j+worksize{
		if j > amount {
			j = amount
		}

		wg.Add(1)
		go func(i, j int) {

			sub := make(map[string]int)

			for w := i; w < j; w++ {
				sub[strings.Trim(strings.ToLower(words[w]), ".,!")]++
			}

			freqsCh <- sub

			wg.Done()
		}(i, j)

	}
	wg.Wait()
	close(freqsCh)
	

	for sub := range freqsCh {
		for k, v := range sub {
			freqs[k] += v
		}
	}

	return freqs
}

// Benchmark how long it takes to count word frequencies in text numRuns times.
//
// Return the total time elapsed.
func benchmark(text string, numRuns int) int64 {
	start := time.Now()
	for i := 0; i < numRuns; i++ {
		WordCount(text)
	}
	runtimeMillis := time.Since(start).Nanoseconds() / 1e6

	return runtimeMillis
}

// Print the results of a benchmark
func printResults(runtimeMillis int64, numRuns int) {
	fmt.Printf("amount of runs: %d\n", numRuns)
	fmt.Printf("total time: %d ms\n", runtimeMillis)
	average := float64(runtimeMillis) / float64(numRuns)
	fmt.Printf("average time/run: %.2f ms\n", average)
}

func main() {
	// read in DataFile as a string called data
	data, err := os.ReadFile("loremipsum.txt")
	if err != nil {
		log.Fatal(err)
	}


	WordCount(string(data))
	//fmt.Printf("%#v", WordCount(string(data)))

	numRuns := 100
	runtimeMillis := benchmark(string(data), numRuns)
	printResults(runtimeMillis, numRuns)
}
