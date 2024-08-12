package main

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
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

func GetContainerByID(ctx context.Context, containerID string) (types.Container, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return types.Container{}, err
	}
	containers, err := cli.ContainerList(ctx, container.ListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "id", Value: containerID}),
	})
	if err != nil {
		return types.Container{}, err
	}
	return containers[0], nil
}

// RunContainers runs containers with the specified images
// And returns the container IDs as a slice of strings
func RunContainers(ctx context.Context, images ...string) ([]string, error) {
	containerIDs := make([]string, 0)
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	for _, img := range images {
		// Pull the image
		reader, err := cli.ImagePull(ctx, img, image.PullOptions{})
		if err != nil {
			return nil, err
		}
		defer reader.Close()
		_, err = io.Copy(os.Stdout, reader)
		if err != nil {
			return nil, err
		}

		// Create a container
		createResponse, err := cli.ContainerCreate(ctx, &container.Config{
			Image: img,
		}, nil, nil, nil, img)
		if err != nil {
			return nil, err
		}
		slog.Info("Created container", "ID", createResponse.ID, "Name", img)

		// Start the container
		err = cli.ContainerStart(ctx, createResponse.ID, container.StartOptions{})
		if err != nil {
			return nil, err
		}
		slog.Info("Started container", "ID", createResponse.ID, "Name", img)
		containerIDs = append(containerIDs, createResponse.ID)
	}
	return containerIDs, nil
}

// StopAndRemoveContainers stops and removes containers with the specified container IDs
func StopAndRemoveContainers(ctx context.Context, containerIDs ...string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	for _, containerID := range containerIDs {
		err = cli.ContainerStop(ctx, containerID, container.StopOptions{})
		if err != nil {
			return err
		}
		slog.Info("Stopped container", "ID", containerID)
		err = cli.ContainerRemove(ctx, containerID, container.RemoveOptions{})
		if err != nil {
			return err
		}
		slog.Info("Removed container", "ID", containerID)
	}
	return nil
}
