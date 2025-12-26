package main

import (
	"fmt"
	"os"

	"github.com/shlokmestry/whyismybuildslow/internal/runner"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "run":
		exitCode, err := runner.Run(os.Args[2:])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
		}
		os.Exit(exitCode)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(2)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  whyismybuildslow run -- <command> [args...]")
}
