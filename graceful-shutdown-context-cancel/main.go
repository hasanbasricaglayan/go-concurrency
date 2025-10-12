package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// worker simulates a long-running task that listens for the cancellation signal
func worker(ctx context.Context, wg *sync.WaitGroup, id int) {
	defer wg.Done()

	fmt.Printf("Worker %d started\n", id)

	// Infinite for loop to simulate continuous work and handle cancellation signal
	for {
		select {
		case <-ctx.Done():
			// Handle context cancellation (graceful shutdown)
			fmt.Printf("Worker %d stopping: %v\n", id, ctx.Err())
			return
		default:
			// Simulate doing some work
			fmt.Printf("Worker %d is working...\n", id)
			time.Sleep(2 * time.Second)
		}
	}
}

func main() {
	// Create a context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())

	// Ensure cancel is called at the end to clean up
	defer cancel()

	// Wait group to ensure all goroutines finish before exiting
	var wg sync.WaitGroup

	// Channel to listen for system signals (e.g., Ctrl+C)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Start multiple worker goroutines
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(ctx, &wg, i)
	}

	// Wait for an interrupt signal to initiate graceful shutdown
	<-signalCh

	// Handle shutdown signal (Ctrl+C or SIGTERM)
	fmt.Println("Received shutdown signal. Shutting down gracefully...")

	// Cancel the context to notify all goroutines to stop
	cancel()

	// Wait for all goroutines to finish
	wg.Wait()

	// Final cleanup before exiting
	fmt.Println("Application shutdown complete.")
}
