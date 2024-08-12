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

func RunContainers(images ...string) error {
	context := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	for _, image := range images {
		// Create a container
		createResponse, err := cli.ContainerCreate(context, &container.Config{
			Image: image,
		}, nil, nil, nil, image)
		if err != nil {
			return err
		}

		// Start the container
		err = cli.ContainerStart(context, createResponse.ID, container.StartOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}
