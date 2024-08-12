package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var ContainersToRun = []string{"nginx"}

func main() {
	containerIDs, err := RunContainers("nginx")
	if err != nil {
		slog.Error("Failed to run container", "error", err)
	}
	Policies, err := GeneratePoliciesForContainers(containerIDs...)
	if err != nil {
		slog.Error("Failed to generate policies", "error", err)
	}
	fmt.Println("Policies generated", Policies)

	// Graceful shutdown
	go func ()  {
		s := make(chan os.Signal, 1)
		signal.Notify(s, os.Interrupt)
		signal.Notify(s, syscall.SIGTERM)
		<-s
		fmt.Println("Shutting down...")
		err := StopAndRemoveContainers(containerIDs...)
		if err != nil {
			slog.Error("Failed to stop and remove containers", "error", err)
		}
		os.Exit(0)
	} ()
	// Block forever
	forever := make(chan bool)
	<-forever
}
