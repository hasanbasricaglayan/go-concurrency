package main

import (
	"fmt"
)

// chan<- string is a send-only channel
func greet(ch chan<- string) {
	// Send message to the channel
	ch <- "Hello from goroutine!"

	// Close the channel only from sender
	close(ch)
}

func main() {
	// Create an unbuffered channel of type string
	messageChannel := make(chan string)

	// Start a goroutine and pass the channel to it
	go greet(messageChannel)

	// Receive the message from the channel
	message := <-messageChannel
	fmt.Println("Received:", message)

	// Check whether the channel is closed
	// ok is false if the channel is closed
	_, ok := <-messageChannel
	if !ok {
		fmt.Println("Channel is closed!")
	}
}
