package main

import (
	"context"

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
