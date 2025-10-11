package main

func main() {
	worker := NewWorker()
	go worker.SendSignal()

	// Wait for the worker to send a signal
	<-worker.signalCh
}
