package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	PROGRAM_TIMEOUT            = 10
	WORKER_COUNT               = 3
	WORKER_PERIOD              = 2
	WORKER_FAILURE_PROBABILITY = 0.01
)

// logger prints structured log with timestamp, level and worker ID
func logger(level string, workerID int, message string) {
	timestamp := time.Now().Format(time.RFC3339)
	if workerID > 0 {
		fmt.Printf("%s [%s] [Worker %d] %s\n", timestamp, level, workerID, message)
	} else {
		fmt.Printf("%s [%s] [Main] %s\n", timestamp, level, message)
	}
}

// worker simulates a long-running periodic task that fails or listens for a cancellation signal
func worker(ctx context.Context, id int, failedWorkerCh chan<- int) error {
	for {
		select {
		case <-time.After(WORKER_PERIOD * time.Second):
			// Simulate doing some work every WORKER_PERIOD seconds
			logger("INFO", id, "Working...")

			// Simulate a random failure
			if rand.Float32() < WORKER_FAILURE_PROBABILITY {
				err := fmt.Errorf("worker failure")
				logger("ERROR", id, err.Error())

				// Send the failing worker ID if channel is empty
				// Default is executed if the ID is already sent
				select {
				case failedWorkerCh <- id:
				default:
				}

				return err
			}

		case <-ctx.Done():
			// Handle context cancellation (graceful shutdown)
			logger("INFO", id, fmt.Sprintf("Stopping: %v", ctx.Err()))

			// Cancellation is treated as normal (not an error)
			return nil
		}
	}
}

func main() {
	// Cancel automatically after PROGRAM_TIMEOUT seconds if no OS signal or worker failure
	ctx, cancel := context.WithTimeout(context.Background(), PROGRAM_TIMEOUT*time.Second)

	// Ensure cancel is called at the end to clean up
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)

	// Channel to capture the ID of the first failing worker
	// Buffered channel to avoid blocking
	failedWorkerCh := make(chan int, 1)

	// Start WORKER_COUNT periodic workers
	for i := 1; i <= WORKER_COUNT; i++ {
		group.Go(func() error {
			return worker(ctx, i, failedWorkerCh)
		})
	}

	// Channel to receive OS signals
	signalCh := make(chan os.Signal, 1)

	// Subscribe to interrupt and terminate signals
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	logger("INFO", 0, "Program running. Press Ctrl+C to stop...")

	// Wait for either OS signal, worker failure or context timeout
	var shutdownReason string
	select {
	case signal := <-signalCh:
		// Cancel the context
		shutdownReason = fmt.Sprintf("Received OS signal: %v", signal)
		cancel()
	case <-ctx.Done():
		select {
		case failedWorkerID := <-failedWorkerCh:
			// Context will cancel automatically on errgroup error (worker failure)
			shutdownReason = fmt.Sprintf("Worker %d caused context cancellation", failedWorkerID)
		default:
			// Context will cancel automatically on timeout
			shutdownReason = "Timeout reached"
		}
	}

	// Wait for all worker goroutines to finish
	err := group.Wait()

	// Shutdown logging
	if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		logger("ERROR", 0, fmt.Sprintf("Shutdown with error: %v | Reason: %s", err, shutdownReason))
	} else {
		logger("INFO", 0, fmt.Sprintf("Shutdown cleanly | Reason: %s", shutdownReason))
	}
}
