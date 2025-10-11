package main

import (
	"fmt"
	"time"
)

func main() {
	// Create a channel to receive a string
	ch := make(chan string)

	// Start a goroutine that waits for 2 seconds and then sends a message
	go func() {
		// Simulate a delay
		time.Sleep(10 * time.Second)

		// Send message after delay
		ch <- "Finally got data!"
	}()

	// Use select to either receive from the channel or timeout
	select {
	case message := <-ch:
		// If data is received from the channel before timeout
		fmt.Println("Received:", message)
	case <-time.After(1 * time.Second):
		// If no data arrives in 1 second, this case runs
		fmt.Println("Timeout! Moving on...")
	}
}
