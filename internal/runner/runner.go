package runner

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/shlokmestry/whyismybuildslow/internal/classifier"
	"github.com/shlokmestry/whyismybuildslow/internal/events"
	"github.com/shlokmestry/whyismybuildslow/internal/ui"
)

func Run(args []string, noUI bool) (int, error) {
	// -----------------------------
	// Argument validation
	// -----------------------------
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

	// -----------------------------
	// Start UI (optional)
	// -----------------------------
	var p *tea.Program

	if !noUI {
		p = tea.NewProgram(ui.InitialModel())
		go func() {
			_ = p.Start()
		}()
	}

	// -----------------------------
	// Event recorder
	// -----------------------------
	recorder := events.NewRecorder()
	recorder.Record("start", "build started")

	start := time.Now()

	// -----------------------------
	// Prepare command
	// -----------------------------
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

	// -----------------------------
	// Start command
	// -----------------------------
	if err := cmd.Start(); err != nil {
		return 1, err
	}

	// Capture output silently (UI owns stdout)
	go scanOutput(stdoutPipe, recorder)
	go scanOutput(stderrPipe, recorder)

	// -----------------------------
	// Wait for completion
	// -----------------------------
	err = cmd.Wait()

	end := time.Now()
	elapsed := end.Sub(start)

	recorder.Record("end", "build finished")

	// -----------------------------
	// Analyze idle gaps
	// -----------------------------
	detectIdleGaps(recorder.Events, 2*time.Second, p)

	// -----------------------------
	// Finish UI cleanly
	// -----------------------------
	if p != nil {
		p.Send(ui.FinishMsg{})
	}

	// -----------------------------
	// Summary output
	// -----------------------------
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

// ------------------------------------------------------
// Output scanning (no printing; UI owns terminal)
// ------------------------------------------------------
func scanOutput(reader io.Reader, recorder *events.Recorder) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		recorder.Record("output", scanner.Text())
	}
}

// ------------------------------------------------------
// Idle gap detection + classification + UI signaling
// ------------------------------------------------------
func detectIdleGaps(
	eventsList []events.Event,
	threshold time.Duration,
	p *tea.Program,
) {
	for i := 1; i < len(eventsList); i++ {
		prev := eventsList[i-1]
		curr := eventsList[i]
		gap := curr.Time.Sub(prev.Time)

		if gap > threshold {
			result := classifier.ClassifyIdleGap(
				prev.Message,
				curr.Message,
				gap.Seconds(),
			)

			// Notify UI only if enabled
			if p != nil {
				p.Send(ui.StallMsg{
					Duration: gap,
					Cause:    string(result.Cause),
				})
			}
		}
	}
}
