package events

import "time"

// Event represents something that happened during the build
type Event struct {
	Time    time.Time
	Type    string // "start", "output", "end"
	Message string
}

// Recorder stores build events in order
type Recorder struct {
	Events []Event
}

// NewRecorder creates a new Recorder
func NewRecorder() *Recorder {
	return &Recorder{
		Events: []Event{},
	}
}

// Record adds a new event with the current timestamp
func (r *Recorder) Record(eventType, message string) {
	r.Events = append(r.Events, Event{
		Time:    time.Now(),
		Type:    eventType,
		Message: message,
	})
}
