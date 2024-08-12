package main

import (
	"context"
	"log/slog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func GetContainers() []types.Container {
	context := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		slog.Error("Failed to initialize new Client")
	}
	containers, err := cli.ContainerList(context, container.ListOptions{})
	if err != nil {
		slog.Error("Failed to list containers")
	}
	return containers
}
