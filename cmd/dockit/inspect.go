package main

import (
	"context"
	"fmt"
	"io"

	imageoptions "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func printImageInfo(image imageoptions.InspectResponse) {
	fmt.Println("ID:", image.ID)
	fmt.Println("Tags:", image.RepoTags)
	fmt.Println("Digests:", image.RepoDigests)
	fmt.Println("Created:", image.Created)
	fmt.Println("Size:", image.Size, "bytes")
	fmt.Println("Architecture:", image.Architecture)
	fmt.Println("OS:", image.Os)
	fmt.Println("Author:", image.Author)

	if image.Config == nil {
		return
	}

	fmt.Println("Entrypoint:", image.Config.Entrypoint)
	fmt.Println("Command:", image.Config.Cmd)
	fmt.Println("Working directory:", image.Config.WorkingDir)
	fmt.Println("User:", image.Config.User)
	fmt.Println("Environment:", image.Config.Env)
	fmt.Println("Labels:", image.Config.Labels)
	fmt.Println("Exposed ports:", image.Config.ExposedPorts)
	fmt.Println("Volumes:", image.Config.Volumes)
	fmt.Println("Layers:", image.RootFS.Layers)
}

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

	pulled := false
	image, err := cli.ImageInspect(ctx, imageName)
	if err != nil {
		fmt.Println("[-] Image not found locally, pulling:", imageName)

		pull, pullErr := cli.ImagePull(ctx, imageName, imageoptions.PullOptions{})
		if pullErr != nil {
			return fmt.Errorf("[-] pull image %q: %w", imageName, pullErr)
		}
		defer pull.Close()

		if _, err := io.Copy(io.Discard, pull); err != nil {
			return fmt.Errorf("[-] download image %q: %w", imageName, err)
		}

		pulled = true
		image, err = cli.ImageInspect(ctx, imageName)
		if err != nil {
			return fmt.Errorf("[-] inspect pulled image %q: %w", imageName, err)
		}
	}

	if pulled {
		fmt.Println("Image pulled:", imageName)
	} else {
		fmt.Println("Image found locally:", imageName)
	}
	printImageInfo(image)

	return nil
}
