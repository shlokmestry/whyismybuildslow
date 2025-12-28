package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
			jsonOut = true
		case "--help", "-h":
			printHelp()
			os.Exit(0)
		default:
			filtered = append(filtered, a)
		}
	}

	// Require at least something
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

	// JSON mode implies no UI
	if jsonOut {
		noUI = true
	}

	exitCode, err := runner.Run(filtered, noUI)
	if err != nil && !jsonOut {
		fmt.Fprintln(os.Stderr, "error:", err)
	}

	// Emit JSON if requested
	if jsonOut {
		summary := map[string]interface{}{
			"command":   strings.Join(filtered[1:], " "),
			"exit_code": exitCode,
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(summary)
	}

	os.Exit(exitCode)
}

func printHelp() {
	fmt.Println(`WhyIsMyBuildSlow 🐌
Usage:
  whyismybuildslow run [flags] -- <command>

Flags:
  --no-ui     Disable animated UI (CI / logs only)
  --json      Output machine-readable JSON
  -h, --help  Show this help

Examples:
  whyismybuildslow run -- npm install
  whyismybuildslow run --no-ui -- sleep 4
  whyismybuildslow run --json -- sleep 2
`)
}
