package main

import "log/slog"

func main() {
	err := RunContainers("nginx")
	if err != nil {
		slog.Error("Failed to run container", "error", err)
	}
}
