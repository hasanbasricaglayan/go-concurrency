package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

const (
	TASK_COUNT = 20
	URL        = "https://jsonplaceholder.typicode.com/posts"
)

func task(wg *sync.WaitGroup, id int) {
	// 3. Tell WaitGroup this task is done = Decrease WaitGroup counter when done
	defer wg.Done()

	resp, err := http.Get(URL)
	if err != nil {
		log.Printf("Task %d failed: %v\n", id, err)
		return
	}

	// Close the response body to avoid memory leak
	defer resp.Body.Close()

	fmt.Printf("Task %d, Response Status: %s\n", id, resp.Status)
}

func main() {
	// 1. Create a WaitGroup = Used to wait for all goroutines to finish
	var wg sync.WaitGroup

	// 2. We are launching TASK_COUNT goroutines
	for i := 1; i <= TASK_COUNT; i++ {
		// Increase WaitGroup counter
		wg.Add(1)

		go task(&wg, i)
	}

	// 4. Wait for all tasks to finish
	wg.Wait()
	fmt.Println("All tasks completed")

	fmt.Println("Main function done")
}
