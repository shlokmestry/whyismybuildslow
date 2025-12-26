package runner

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/shlokmestry/whyismybuildslow/internal/events"
)

func Run(args []string) (int, error) {
	// --- Validate args ---
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

	// --- Event recorder ---
	recorder := events.NewRecorder()
	recorder.Record("start", "build started")

	start := time.Now()

	fmt.Printf("🐌 WhyIsMyBuildSlow starting at %s\n", start.Format(time.RFC3339))
	fmt.Printf("🐌 Running: %s %s\n", command, join(commandArgs))

	// --- Prepare command ---
	cmd := exec.Command(command, commandArgs...)
	cmd.Stdin = os.Stdin

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return 1, err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return 1, err
	}

	// --- Start command ---
	if err := cmd.Start(); err != nil {
		return 1, err
	}

	// --- Capture output ---
	go scanOutput(stdoutPipe, recorder)
	go scanOutput(stderrPipe, recorder)

	// --- Wait for command ---
	err = cmd.Wait()

	end := time.Now()
	elapsed := end.Sub(start)

	recorder.Record("end", "build finished")

	fmt.Printf("\n🐌 Finished at %s\n", end.Format(time.RFC3339))
	fmt.Printf("🐌 Elapsed time: %s\n", elapsed)

	// --- Detect idle gaps ---
	detectIdleGaps(recorder.Events, 2*time.Second)

	// --- Exit code handling ---
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

// scanOutput reads output line-by-line and records events
func scanOutput(reader io.Reader, recorder *events.Recorder) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		recorder.Record("output", line)
	}
}

// detectIdleGaps finds periods where nothing happened
func detectIdleGaps(eventsList []events.Event, threshold time.Duration) {
	for i := 1; i < len(eventsList); i++ {
		prev := eventsList[i-1]
		curr := eventsList[i]
		gap := curr.Time.Sub(prev.Time)

		if gap > threshold {
			fmt.Printf(
				"\n🐌 Idle gap detected: %s\n",
				gap.Round(time.Millisecond),
			)
		}
	}
}

// join helper for printing command args
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
