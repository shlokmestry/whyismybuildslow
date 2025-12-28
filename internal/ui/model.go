package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Messages
type TickMsg time.Time
type FinishMsg struct{}
type StallMsg struct {
	Duration time.Duration
	Cause    string
}

// Model
type Model struct {
	StartTime time.Time
	Now       time.Time
	Finished  bool

	Stalled   bool
	StallFor  time.Duration
	Cause     string
}

// Initial model
func InitialModel() Model {
	now := time.Now()
	return Model{
		StartTime: now,
		Now:       now,
	}
}

// Init
func (m Model) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Update
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case TickMsg:
		m.Now = time.Time(msg)
		return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
			return TickMsg(t)
		})

	case StallMsg:
		m.Stalled = true
		m.StallFor = msg.Duration
		m.Cause = msg.Cause
		return m, nil

	case FinishMsg:
		m.Finished = true
		return m, tea.Quit
	}

	return m, nil
}

// View
func (m Model) View() string {
	elapsed := m.Now.Sub(m.StartTime)
	bar := progressBar(elapsed)

	if m.Finished {
		return fmt.Sprintf(
			"✅ Build finished\nElapsed: %s\n",
			elapsed.Round(time.Millisecond),
		)
	}

	if m.Stalled {
		return fmt.Sprintf(
			"%s\nElapsed: %s\n\n%s\n",
			bar,
			elapsed.Round(time.Millisecond),
			characterLine(m.Cause),
		)
	}

	return fmt.Sprintf(
		"%s\nElapsed: %s\n",
		bar,
		elapsed.Round(time.Millisecond),
	)
}



func progressBar(elapsed time.Duration) string {
	// Fake progress for UX (we don't know total time)
	percent := int(elapsed.Seconds()*15) % 100
	width := 20
	filled := (percent * width) / 100

	return fmt.Sprintf(
		"⚙️ Building %s%s %d%%",
		strings.Repeat("█", filled),
		strings.Repeat("░", width-filled),
		percent,
	)
}

func characterLine(cause string) string {
	switch cause {
	case "network":
		return "🐌 Network Slug crawling…"
	case "cache":
		return "🧊 Cache Golem awakening…"
	case "docker":
		return "🚚 Docker Truck unloading layers…"
	default:
		return "🤷 Something is slowing things down…"
	}
}
