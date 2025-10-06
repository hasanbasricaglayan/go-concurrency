package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

func backgroundTask(id int, wg *sync.WaitGroup) {
	// 3. Tell WaitGroup this task is done = Decrease WaitGroup counter when done
	defer wg.Done()

	url := "https://jsonplaceholder.typicode.com/posts"
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Task %d failed: %v\n", id, err)
		return
	}

	fmt.Printf("Task %d, Response Status: %s\n", id, resp.Status)
}

func main() {
	// 1. Create a WaitGroup = Used to wait for all goroutines to finish
	var wg sync.WaitGroup

	// 2. We are launching 20 goroutines
	totalTasks := 20
	for i := 1; i <= totalTasks; i++ {
		// Increase WaitGroup counter
		wg.Add(1)

		go backgroundTask(i, &wg)
	}

	// 4. Wait for all tasks to finish
	wg.Wait()
	fmt.Println("All background tasks completed")

	fmt.Println("Main function done")
}
