package main

import (
	"fmt"
	"sync"
)

// chan<- string is a send-only channel
func worker(id int, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	ch <- fmt.Sprintf("Worker %d done", id)
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan string)

	numWorkers := 5
	wg.Add(numWorkers)

	// Start multiple sender goroutines
	for i := 1; i <= numWorkers; i++ {
		go worker(i, ch, &wg)
	}

	// Dedicated goroutine to close the channel
	go func() {
		// Wait until all workers are done
		wg.Wait()

		// Only one closer
		close(ch)
	}()

	// Receive from the channel until it's closed
	for message := range ch {
		fmt.Println("Received:", message)
	}
}
