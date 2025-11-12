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
	PROGRAM_TIMEOUT             = 11
	GENERATOR_MAX_INT           = 10
	GENERATOR_DURATION          = 1
	GENERATOR_ERROR_PROBABILITY = 0.01
)

// generator returns a function that emits numbers from 1 to max on dataCh
func generator(ctx context.Context, max int, dataCh chan<- int) func() error {
	return func() error {
		defer close(dataCh)

		for i := 1; i <= max; i++ {
			// Wait for either the duration or context cancellation
			select {
			case <-time.After(GENERATOR_DURATION * time.Second):

			case <-ctx.Done():
				// Context canceled
				return ctx.Err()
			}

			// Simulate a random error
			if rand.Float32() < GENERATOR_ERROR_PROBABILITY {
				return errors.New("generator error")
			}

			// Send number to the channel, respecting context cancellation
			select {
			case dataCh <- i:

			case <-ctx.Done():
				// Context canceled
				return ctx.Err()
			}
		}

		return nil
	}
}

// receiver processes data from dataCh and handles OS signals
func receiver(ctx context.Context, cancel context.CancelFunc, dataCh <-chan int, signalCh <-chan os.Signal) func() error {
	return func() error {
		for {
			select {
			case <-ctx.Done():
				// Context will cancel automatically on timeout
				return ctx.Err()

			case sig := <-signalCh:
				// Cancel the parent context (ctx)
				cancel()
				return fmt.Errorf("received OS signal: %v", sig)

			case v, ok := <-dataCh:
				if !ok {
					// Generator finished successfully
					return nil
				}
				fmt.Printf("Received: %d\n", v)
			}
		}
	}
}

func main() {
	// Create a cancellable context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), PROGRAM_TIMEOUT*time.Second)

	// Ensure resources are cleaned up
	defer cancel()

	// Channel to receive OS signals
	signalCh := make(chan os.Signal, 1)

	// Subscribe to interrupt and terminate signals
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// groupCtx is derived from parent ctx
	group, groupCtx := errgroup.WithContext(ctx)

	dataCh := make(chan int)

	group.Go(generator(groupCtx, GENERATOR_MAX_INT, dataCh))
	group.Go(receiver(groupCtx, cancel, dataCh, signalCh))

	// Wait for the goroutines or first error
	if err := group.Wait(); err != nil {
		fmt.Println("Shutdown reason:", err)
	} else {
		fmt.Println("Completed successfully")
	}
}
