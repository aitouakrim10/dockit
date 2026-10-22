package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Dockit is starting...")
	if len(os.Args) > 1 {
		if os.Args[1] == "-h" || os.Args[1] == "help" {
			fmt.Println("Help information for Dockit")
			fmt.Println("Usage: dockit [options]")
			fmt.Println("Options:")
			fmt.Println("  -h, help     Show this help message and exit")
			fmt.Println("  inspect  Inspect a Docker image")
			fmt.Println("  scan     Scan a Docker image for vulnerabilities")
		} else if os.Args[1] == "inspect" && len(os.Args) == 3 {
			fmt.Println("Inspecting Docker image...", os.Args[2])
		} else if os.Args[1] == "scan" && len(os.Args) == 3 {
			fmt.Println("Scanning Docker image for vulnerabilities...", os.Args[2])
		}
	} else {
		fmt.Println("No arguments provided. Use -h or --help for usage information.")
	}
	os.Exit(0)
}
