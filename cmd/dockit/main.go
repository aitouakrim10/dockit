package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Dockit is starting...")
	if len(os.Args) > 1 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" || os.Args[1] == "help" {
			fmt.Println("Help information for Dockit")
			fmt.Println("Usage: dockit <command> [arguments]")
			fmt.Println("Commands:")
			fmt.Println("  inspect <image>  Inspect a Docker image")
			fmt.Println("  scan <image>     Scan a Docker image for vulnerabilities")
		} else if os.Args[1] == "inspect" && len(os.Args) == 3 {
			if err := inspectImage(os.Args[2]); err != nil {
				fmt.Println("Error:", err)
			}
		} else if os.Args[1] == "scan" && len(os.Args) == 3 {
			fmt.Println("Scanning Docker image for vulnerabilities...", os.Args[2])
		}
	} else {
		fmt.Println("No arguments provided. Use -h or --help for usage information.")
	}
	os.Exit(0)
}
