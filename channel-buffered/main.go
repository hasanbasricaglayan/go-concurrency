package main

import (
	"fmt"
)

func main() {
	// Create a buffered string channel of size 2
	ch := make(chan string, 2)

	// Send messages to the channel
	ch <- "One"
	ch <- "Two"

	// Need another goroutine to send the third value as the buffer is full
	go func() {
		ch <- "Three"
		close(ch)
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
