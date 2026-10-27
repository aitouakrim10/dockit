package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" || os.Args[1] == "help" {
			fmt.Println("Help information for Dockit")
			fmt.Println("Usage: dockit <command> [arguments]")
			fmt.Println("Commands:")
			fmt.Println("  inspect <image> [--format text|json|yaml]  Inspect a Docker image")
			fmt.Println("  scan <image>     Scan a Docker image for vulnerabilities")
		} else if os.Args[1] == "inspect" {
			imageName, format, err := parseInspectArgs(os.Args[2:])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				os.Exit(1)
			}
			if err := inspectImage(imageName, format); err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				os.Exit(1)
			}
		} else if os.Args[1] == "scan" && len(os.Args) == 3 {
			if err := scanImage(os.Args[2]); err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				os.Exit(1)
			}
		}
	} else {
		fmt.Println("No arguments provided. Use -h or --help for usage information.")
	}
	os.Exit(0)
}
