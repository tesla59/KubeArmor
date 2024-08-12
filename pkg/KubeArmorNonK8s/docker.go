package main

import (
	"context"
	"log/slog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func GetContainers() ([]types.Container, error) {
	context := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	containers, err := cli.ContainerList(context, container.ListOptions{})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// RunContainers runs containers with the specified images
// And returns the container IDs as a slice of strings
func RunContainers(images ...string) ([]string, error) {
	containerIDs := make([]string, 0)
	context := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	for _, image := range images {
		// Create a container
		createResponse, err := cli.ContainerCreate(context, &container.Config{
			Image: image,
		}, nil, nil, nil, image)
		if err != nil {
			return nil, err
		}
		slog.Info("Created container", "ID", createResponse.ID, "Name", image)

		// Start the container
		err = cli.ContainerStart(context, createResponse.ID, container.StartOptions{})
		if err != nil {
			return nil, err
		}
		slog.Info("Started container", "ID", createResponse.ID, "Name", image)
		containerIDs = append(containerIDs, createResponse.ID)
	}
	return containerIDs, nil
}

// StopAndRemoveContainers stops and removes containers with the specified container IDs
func StopAndRemoveContainers(containerIDs ...string) error {
	context := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	for _, containerID := range containerIDs {
		err = cli.ContainerStop(context, containerID, container.StopOptions{})
		if err != nil {
			return err
		}
		slog.Info("Stopped container", "ID", containerID)
	}
	return nil
}
