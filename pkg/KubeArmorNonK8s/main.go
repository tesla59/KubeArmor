package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	_, err := RunContainers("nginx")
	if err != nil {
		slog.Error("Failed to run container", "error", err)
	}

	// Graceful shutdown
	go func ()  {
		s := make(chan os.Signal, 1)
		signal.Notify(s, os.Interrupt)
		signal.Notify(s, syscall.SIGTERM)
		<-s
		fmt.Println("Shutting down...")
		os.Exit(0)
	} ()
	// Block forever
	forever := make(chan bool)
	<-forever
}
