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

	// Mutex to protect the counter
	var mu sync.Mutex

	// for i := 0; i < 1000; i++
	for range 1000 {
		// Increase WaitGroup counter
		wg.Add(1)

		go func() {
			// Decrease WaitGroup counter when done
			defer wg.Done()

			// Lock before accessing counter
			mu.Lock()

			// Unlock after done
			defer mu.Unlock()

			counter++
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Always 1000
	fmt.Println("Final counter:", counter)
}
