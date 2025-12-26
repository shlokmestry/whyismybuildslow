package runner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func Run(args []string) (int, error) {
	if len(args) == 0 {
		return 2, errors.New("missing '-- <command>'")
	}

	if args[0] != "--" {
		return 2, errors.New("expected '--' before the command")
	}

	if len(args) < 2 {
		return 2, errors.New("missing command after '--'")
	}

	command := args[1]
	commandArgs := []string{}
	if len(args) > 2 {
		commandArgs = args[2:]
	}

	start := time.Now()

	fmt.Printf("🐌 WhyIsMyBuildSlow starting at %s\n", start.Format(time.RFC3339))
	fmt.Printf("🐌 Running: %s %s\n", command, join(commandArgs))

	cmd := exec.Command(command, commandArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()

	end := time.Now()
	elapsed := end.Sub(start)

	fmt.Printf("\n🐌 Finished at %s\n", end.Format(time.RFC3339))
	fmt.Printf("🐌 Elapsed time: %s\n", elapsed)

	if err == nil {
		fmt.Println("🐌 Exit code: 0")
		return 0, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		code := exitErr.ExitCode()
		fmt.Printf("🐌 Exit code: %d\n", code)
		return code, err
	}

	return 1, err
}

func join(args []string) string {
	if len(args) == 0 {
		return ""
	}
	out := ""
	for i, a := range args {
		if i > 0 {
			out += " "
		}
		out += a
	}
	return out
}
