package main

import (
	"fmt"
	"sync"
)

func main() {
	// Shared variable
	counter := 0

	// Used to wait for all goroutines to finish
	var wg sync.WaitGroup

	// for i := 0; i < 1000; i++
	for range 1000 {
		// Increase WaitGroup counter
		wg.Add(1)

		go func() {
			// Decrease WaitGroup counter when done
			defer wg.Done()

			// Race condition happens here
			counter++
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Unpredictable result
	fmt.Println("Final counter:", counter)
}
