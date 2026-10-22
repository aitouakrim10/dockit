package main

import (
	"context"
	"fmt"
	"io"

	imageoptions "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func inspectImage(imageName string) error {
	fmt.Println("Inspecting Docker image:", imageName)
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return err
	}
	defer cli.Close()

	image, err := cli.ImageInspect(ctx, imageName)
	if err != nil {
		fmt.Println("Image not found locally, pulling:", imageName)

		pull, pullErr := cli.ImagePull(ctx, imageName, imageoptions.PullOptions{})
		if pullErr != nil {
			return fmt.Errorf("pull image %q: %w", imageName, pullErr)
		}
		defer pull.Close()

		if _, err := io.Copy(io.Discard, pull); err != nil {
			return fmt.Errorf("download image %q: %w", imageName, err)
		}

		image, err = cli.ImageInspect(ctx, imageName)
		if err != nil {
			return fmt.Errorf("inspect pulled image %q: %w", imageName, err)
		}
	}

	fmt.Println("Image found locally:", imageName)
	fmt.Println("ID:", image.ID)

	return nil
}
