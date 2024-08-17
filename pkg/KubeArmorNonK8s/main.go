package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// ContainersToRun is a list of container images to run
var ContainersToRun = []string{"nginx", "alpine"}

// PoliciesDirectory is the directory where the generated policies are stored
var PoliciesDirectory = "policies"

func main() {
	ctx := context.Background()

	// 1. Run containers
	containerIDs, err := RunContainers(ctx, ContainersToRun...)
	if err != nil {
		slog.Error("Failed to run container", "error", err)
	}

	// 2. Generate policies
	generatedPolicies, err := GeneratePoliciesForContainers(ctx, containerIDs...)
	if err != nil {
		slog.Error("Failed to generate policies", "error", err)
	}
	fmt.Println("Policies generated", generatedPolicies)

	// 3. Apply policies
	err = ApplyPolicies(ctx, PoliciesDirectory)
	if err != nil {
		slog.Info("Failed to apply policy", "error", err)
	}

	// 4. Handle shutdown
	go func() {
		s := make(chan os.Signal, 1)
		signal.Notify(s, os.Interrupt)
		signal.Notify(s, syscall.SIGTERM)
		<-s
		fmt.Println("Shutting down...")
		err := StopAndRemoveContainers(ctx, containerIDs...)
		if err != nil {
			slog.Error("Failed to stop and remove containers", "error", err)
		}
		os.Exit(0)
	}()
	// Block forever
	forever := make(chan bool)
	<-forever
}
