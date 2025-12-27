package classifier
import "strings"


type Cause string

const (
	CauseUnknown   Cause = "unknown"
	CauseNetwork   Cause = "network"
	CauseCache     Cause = "cache"
	CauseDocker    Cause = "docker"
)

type Result struct {
	Cause       Cause
	Character   string
	Explanation string
	Confidence  float64
}

func ClassifyIdleGap(
	prevOutput string,
	nextOutput string,
	gapSeconds float64,
) Result {

	// Very simple heuristics (for now)
	if looksLikeNetwork(prevOutput) || looksLikeNetwork(nextOutput) {
		return Result{
			Cause:       CauseNetwork,
			Character:   "🐌 Network Slug",
			Explanation: "Build was waiting on network activity (likely downloading dependencies)",
			Confidence:  0.85,
		}
	}

	return Result{
		Cause:       CauseUnknown,
		Character:   "🤷 Unknown",
		Explanation: "Could not confidently determine the cause",
		Confidence:  0.3,
	}
}

func looksLikeNetwork(output string) bool {
	networkKeywords := []string{
		"download",
		"fetch",
		"http",
		"https",
		"registry",
		"pull",
		"resolved",
		"connecting",
		"request",
	}

	outputLower := strings.ToLower(output)

	for _, keyword := range networkKeywords {
		if strings.Contains(outputLower, keyword) {
			return true
		}
	}

	return false
}

