package main

import (
	"fmt"
	"time"
)

func main() {
	fast := make(chan string)
	slow := make(chan string)

	// The first goroutine sleeps for 1 second and sends "I'm fast!"
	go func() {
		time.Sleep(1 * time.Second)
		fast <- "I'm fast!"
	}()

	// The second goroutine sleeps for 2 seconds and sends "I'm slow!"
	go func() {
		time.Sleep(2 * time.Second)
		slow <- "I'm slow!"
	}()

	// select waits until any one of the channels is ready
	// Because fast sends first, that case runs, and we skip the slower one
	select {
	case message := <-fast:
		fmt.Println("Got:", message)
	case message := <-slow:
		fmt.Println("Got:", message)
	}
}
