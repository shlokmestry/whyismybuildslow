package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Messages sent to the UI
type TickMsg time.Time
type StartMsg struct{}
type FinishMsg struct{}
type StallMsg struct {
	Duration time.Duration
}

// UI model
type Model struct {
	StartTime time.Time
	Now       time.Time
	Stalled   bool
	StallFor  time.Duration
	Finished  bool
}

// Create initial model
func InitialModel() Model {
	return Model{
		StartTime: time.Now(),
		Now:       time.Now(),
	}
}

// Init is called once
func (m Model) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Update handles messages
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
		return m, nil

	case FinishMsg:
		m.Finished = true
		return m, tea.Quit
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	elapsed := m.Now.Sub(m.StartTime).Round(time.Millisecond)

	if m.Finished {
		return "✅ Build finished\n"
	}

	if m.Stalled {
		return "🐌 Network Slug crawling...\nElapsed: " + elapsed.String() + "\n"
	}

	return "⚙️ Building...\nElapsed: " + elapsed.String() + "\n"
}
