package main

import (
	"fmt"
	"os"

	"github.com/shlokmestry/whyismybuildslow/internal/runner"
)

func main() {
	args := os.Args[1:]

	noUI := false
	jsonOut := false
	filtered := []string{}

	for _, a := range args {
		switch a {
		case "--no-ui":
			noUI = true
		case "--json":
			jsonOut = true // reserved for Week 7
		case "--help", "-h":
			printHelp()
			os.Exit(0)
		default:
			filtered = append(filtered, a)
		}
	}

	// Guard: require at least something
	if len(filtered) == 0 {
		printHelp()
		os.Exit(2)
	}

	// Strip "run" subcommand if present
	if filtered[0] == "run" {
		filtered = filtered[1:]
	}

	// After stripping, we still need the actual command
	if len(filtered) == 0 {
		printHelp()
		os.Exit(2)
	}

	_ = jsonOut // intentionally unused for now

	code, err := runner.Run(filtered, noUI)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
	}

	os.Exit(code)
}

func printHelp() {
	fmt.Println(`WhyIsMyBuildSlow 🐌
Usage:
  whyismybuildslow run [flags] -- <command>

Flags:
  --no-ui     Disable animated UI (CI / logs only)
  --json      Output machine-readable JSON (coming soon)
  -h, --help  Show this help

Examples:
  whyismybuildslow run -- npm install
  whyismybuildslow run --no-ui -- sleep 4
`)
}
