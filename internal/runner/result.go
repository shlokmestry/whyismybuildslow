package runner

// Summary is the machine-readable result of a build run
type Summary struct {
	Command  string  `json:"command"`
	Elapsed  float64 `json:"elapsed_seconds"`
	Cause    string  `json:"cause"`
	ExitCode int     `json:"exit_code"`
}
