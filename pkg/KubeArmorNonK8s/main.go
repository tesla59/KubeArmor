package main

import "log/slog"

func main() {
	_, err := RunContainers("nginx")
	if err != nil {
		slog.Error("Failed to run container", "error", err)
	}
}
