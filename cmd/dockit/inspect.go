package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	imageoptions "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"go.yaml.in/yaml/v3"
)

type outputFormat string

const (
	textFormat outputFormat = "text"
	jsonFormat outputFormat = "json"
	yamlFormat outputFormat = "yaml"
)

func parseInspectArgs(args []string) (string, outputFormat, error) {
	if len(args) < 1 {
		return "", "", fmt.Errorf("usage: dockit inspect <image> [--format text|json|yaml]")
	}

	format := textFormat
	imageName := ""
	for index := 0; index < len(args); index++ {
		switch args[index] {
		case "--format", "-f":
			if index+1 >= len(args) {
				return "", "", fmt.Errorf("missing value for %s", args[index])
			}
			format = outputFormat(args[index+1])
			index++
		default:
			if imageName != "" {
				return "", "", fmt.Errorf("unexpected argument %q", args[index])
			}
			imageName = args[index]
		}
	}

	switch format {
	case textFormat, jsonFormat, yamlFormat:
		return imageName, format, nil
	default:
		return "", "", fmt.Errorf("unsupported format %q; choose text, json, or yaml", format)
	}
}

func printImageInfo(image imageoptions.InspectResponse, format outputFormat) error {
	switch format {
	case jsonFormat:
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(image)
	case yamlFormat:
		data, err := yaml.Marshal(image)
		if err != nil {
			return fmt.Errorf("marshal image as YAML: %w", err)
		}
		_, err = os.Stdout.Write(data)
		return err
	case textFormat:
		printTextImageInfo(image)
		return nil
	default:
		return fmt.Errorf("unsupported format %q", format)
	}
}

func printTextImageInfo(image imageoptions.InspectResponse) {
	fmt.Println("\nImage metadata:")
	fmt.Println("ID:", image.ID)
	printStrings("Tags", image.RepoTags)
	printStrings("Digests", image.RepoDigests)
	fmt.Println("Created:", image.Created)
	fmt.Println("Size:", image.Size, "bytes")
	fmt.Println("Architecture:", image.Architecture)
	fmt.Println("OS:", image.Os)
	printString("Author", image.Author)

	if image.Config == nil {
		return
	}

	fmt.Println("\nRuntime configuration:")
	printStrings("Entrypoint", image.Config.Entrypoint)
	printStrings("Command", image.Config.Cmd)
	printString("Working directory", image.Config.WorkingDir)
	printString("User", image.Config.User)
	printStrings("Environment", image.Config.Env)
	printMap("Labels", image.Config.Labels)
	printMap("Exposed ports", image.Config.ExposedPorts)
	printMap("Volumes", image.Config.Volumes)
	printStrings("Layers", image.RootFS.Layers)
}

func printString(label, value string) {
	if value == "" {
		fmt.Printf("%s: none\n", label)
		return
	}
	fmt.Printf("%s: %s\n", label, value)
}

func printStrings(label string, values []string) {
	if len(values) == 0 {
		fmt.Printf("%s: none\n", label)
		return
	}
	fmt.Printf("%s: %v\n", label, values)
}

func printMap[K comparable, V any](label string, value map[K]V) {
	if len(value) == 0 {
		fmt.Printf("%s: none\n", label)
		return
	}
	fmt.Printf("%s: %v\n", label, value)
}

func inspectImage(imageName string, format outputFormat) error {
	if format == textFormat {
		fmt.Println("Inspecting Docker image:", imageName)
	}
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
		fmt.Fprintln(os.Stderr, "[-] Image not found locally, pulling:", imageName)

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

	if format == textFormat {
		if pulled {
			fmt.Println("Image pulled:", imageName)
		} else {
			fmt.Println("Image found locally:", imageName)
		}
	} else if pulled {
		fmt.Fprintln(os.Stderr, "Image pulled:", imageName)
	}

	return printImageInfo(image, format)
}
