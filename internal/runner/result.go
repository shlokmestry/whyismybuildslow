package runner

type Summary struct {
	Command   string  `json:"command"`
	Elapsed   float64 `json:"elapsed_seconds"`
	Cause     string  `json:"cause"`
	ExitCode  int     `json:"exit_code"`
}
