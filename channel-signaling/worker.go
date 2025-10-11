package main

import "fmt"

// struct{} type occupies zero bytes of memory
type Worker struct {
	signalCh chan struct{}
}

func NewWorker() *Worker {
	return &Worker{
		signalCh: make(chan struct{}),
	}
}

func (w *Worker) SendSignal() {
	fmt.Println("Signaling!")

	w.signalCh <- struct{}{}
}
