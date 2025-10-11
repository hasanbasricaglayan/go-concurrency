package main

import (
	"fmt"
	"time"
)

func main() {
	fast := make(chan string)
	slow := make(chan string)

	// The first goroutine sleeps for 2 seconds and sends "I'm fast!"
	go func() {
		time.Sleep(2 * time.Second)
		fast <- "I'm fast!"
	}()

	// The second goroutine sleeps for 3 seconds and sends "I'm slow!"
	go func() {
		time.Sleep(3 * time.Second)
		slow <- "I'm slow!"
	}()

	// Try to receive before any goroutine sends
	// default is chosen
	select {
	case message := <-fast:
		fmt.Println("Got:", message)
	case message := <-slow:
		fmt.Println("Got:", message)
	default:
		fmt.Println("No messages yet. Doing something else.")
	}

	// Wait for messages to arrive
	time.Sleep(3 * time.Second)

	// Try again after waiting
	// All channels are ready
	// When more than one case is ready, Go picks one at random
	// Which one runs may change each time the program is ran
	select {
	case message := <-fast:
		fmt.Println("Later got:", message)
	case message := <-slow:
		fmt.Println("Later got:", message)
	default:
		fmt.Println("Still nothing...")
	}
}
