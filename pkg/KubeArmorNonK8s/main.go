package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var ContainersToRun = []string{"nginx"}

func main() {
	ctx := context.Background()

	// 1. Run containers
	containerIDs, err := RunContainers(ContainersToRun...)
	if err != nil {
		slog.Error("Failed to run container", "error", err)
	}

	// 2. Generate policies
	generatedPolicies, err := GeneratePoliciesForContainers(containerIDs...)
	if err != nil {
		slog.Error("Failed to generate policies", "error", err)
	}
	fmt.Println("Policies generated", generatedPolicies)

	// 3. Apply policies
	err = ApplyPolicy(ctx, "policies/policy_nginx.yaml")
	if err != nil {
		slog.Info("Failed to apply policy", "error", err)
	}

	// 4. Handle shutdown
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
