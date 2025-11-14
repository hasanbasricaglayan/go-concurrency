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
	PROGRAM_TIMEOUT            = 12
	WORKER_COUNT               = 3
	WORKER_PERIOD              = 2
	WORKER_FAILURE_PROBABILITY = 0.01
	WORKER_CYCLES              = 5
)

// Custom error type to carry worker ID
type WorkerError struct {
	WorkerID int
	Err      error
}

// In Go, any type that implements the Error() string method automatically satisfies the error interface
func (e *WorkerError) Error() string {
	return fmt.Sprintf("worker %d failed: %v", e.WorkerID, e.Err)
}

func (e *WorkerError) Unwrap() error {
	return e.Err
}

// Signal error type
type SignalError struct {
	Signal os.Signal
}

func (e *SignalError) Error() string {
	return fmt.Sprintf("received OS signal: %v", e.Signal)
}

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
func worker(ctx context.Context, id int) error {
	// Use time.Ticker instead of time.After to avoid creating a new timer on every loop iteration
	ticker := time.NewTicker(WORKER_PERIOD * time.Second)
	defer ticker.Stop()

	// Worker stops after WORKER_CYCLES cycles
	for i := 1; i <= WORKER_CYCLES; i++ {
		select {
		case <-ticker.C:
			// Simulate doing some work every WORKER_PERIOD seconds
			logger("INFO", id, "Working...")

			// Simulate a random failure
			if rand.Float32() < WORKER_FAILURE_PROBABILITY {
				err := errors.New("worker failure")
				logger("ERROR", id, err.Error())

				// WorkerError is an error type
				return &WorkerError{WorkerID: id, Err: err}
			}

		case <-ctx.Done():
			// Handle context cancellation (graceful shutdown)
			logger("INFO", id, fmt.Sprintf("Stopping: %v", ctx.Err()))

			// Worker propagates the cancellation reason
			return ctx.Err()
		}
	}

	// Worker completes its work
	logger("INFO", id, "Done")
	return nil
}

func main() {
	// Cancel automatically after PROGRAM_TIMEOUT seconds if no OS signal or worker failure
	ctx, cancel := context.WithTimeout(context.Background(), PROGRAM_TIMEOUT*time.Second)

	// Ensure cancel is called at the end to clean up
	defer cancel()

	group, groupCtx := errgroup.WithContext(ctx)

	// Start WORKER_COUNT periodic workers
	for i := 1; i <= WORKER_COUNT; i++ {
		// Capture loop variable for Go < 1.22
		i := i
		group.Go(func() error {
			return worker(groupCtx, i)
		})
	}

	// Channel to receive OS signals
	signalCh := make(chan os.Signal, 1)

	// Subscribe to interrupt and terminate signals
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	logger("INFO", 0, "Program running. Press Ctrl+C to stop...")

	// Signal handler goroutine
	group.Go(func() error {
		select {
		case sig := <-signalCh:
			logger("INFO", 0, fmt.Sprintf("Received OS signal: %v", sig))
			cancel()
			return &SignalError{Signal: sig}
		case <-groupCtx.Done():
			return nil
		}
	})

	// Wait for all worker goroutines to finish
	err := group.Wait()

	// Log based on error type
	var workerErr *WorkerError
	var signalErr *SignalError

	switch {
	case err == nil:
		logger("INFO", 0, "Shutdown: completed successfully")
	case errors.As(err, &workerErr):
		logger("ERROR", 0, fmt.Sprintf("Shutdown: %v", err))
	case errors.As(err, &signalErr):
		logger("INFO", 0, fmt.Sprintf("Shutdown: %v", err))
	case errors.Is(err, context.DeadlineExceeded):
		logger("INFO", 0, "Shutdown: timeout reached")
	case errors.Is(err, context.Canceled):
		logger("INFO", 0, "Shutdown: context canceled")
	default:
		logger("ERROR", 0, fmt.Sprintf("Shutdown: unexpected error: %v", err))
	}
}
