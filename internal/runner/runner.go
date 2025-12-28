package runner

import (
	"bufio"
	"errors"
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

	if err := cmd.Start(); err != nil {
		return 1, err
	}

	// Capture output silently (UI owns stdout)
	go scanOutput(stdoutPipe, recorder)
	go scanOutput(stderrPipe, recorder)

	err = cmd.Wait()

	end := time.Now()
	elapsed := end.Sub(start)

	recorder.Record("end", "build finished")

	// -----------------------------
	// Detect stalls
	// -----------------------------
	_ = detectIdleGaps(recorder.Events, 2*time.Second, p)

	if p != nil {
		p.Send(ui.FinishMsg{})
	}

	if err == nil {
		return 0, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode(), err
	}

	return 1, err
}

func scanOutput(reader io.Reader, recorder *events.Recorder) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		recorder.Record("output", scanner.Text())
	}
}

func detectIdleGaps(
	eventsList []events.Event,
	threshold time.Duration,
	p *tea.Program,
) string {

	lastCause := "unknown"

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

			lastCause = string(result.Cause)

			// Notify UI only if enabled
			if p != nil {
				p.Send(ui.StallMsg{
					Duration: gap,
					Cause:    lastCause,
				})
			}
		}
	}

	return lastCause
}
