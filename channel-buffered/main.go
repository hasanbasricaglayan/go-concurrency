package main

import (
	"fmt"
)

func main() {
	// Create a buffered string channel of size 2
	ch := make(chan string, 2)

	// Send 2 messages to the channel
	ch <- "One"
	ch <- "Two"

	// Need another goroutine to send more values as the buffer is full
	go func() {
		// Close the channel to signal no more values will be sent
		defer close(ch)

		ch <- "Three"
		ch <- "Four"
	}()

	// for range ch loop automatically stops when the channel ch is closed
	for message := range ch {
		// Receive the messages from the channel
		fmt.Println(message)
	}

	// Receive a zero value ("" in this case) after the channel is empty
	value, ok := <-ch
	if !ok {
		fmt.Println(value)
	}
}
