package classifier

import "strings"

// Cause represents the type of bottleneck
type Cause string

const (
	CauseUnknown Cause = "unknown"
	CauseNetwork Cause = "network"
	CauseCache   Cause = "cache"
	CauseDocker  Cause = "docker"
)

// Result is the classification output
type Result struct {
	Cause       Cause
	Character   string
	Explanation string
	Confidence  float64
}

// ClassifyIdleGap determines why an idle gap occurred
func ClassifyIdleGap(
	prevOutput string,
	nextOutput string,
	gapSeconds float64,
) Result {

	// Priority matters: Docker > Cache > Network

	if looksLikeDocker(prevOutput) || looksLikeDocker(nextOutput) {
		return Result{
			Cause:       CauseDocker,
			Character:   "🚚 Docker Truck",
			Explanation: "Build was pulling or extracting Docker images",
			Confidence:  0.9,
		}
	}

	if looksLikeCache(prevOutput) || looksLikeCache(nextOutput) {
		return Result{
			Cause:       CauseCache,
			Character:   "🧊 Cache Golem",
			Explanation: "Build cache was cold or invalidated",
			Confidence:  0.8,
		}
	}

	if looksLikeNetwork(prevOutput) || looksLikeNetwork(nextOutput) {
		return Result{
			Cause:       CauseNetwork,
			Character:   "🐌 Network Slug",
			Explanation: "Build was waiting on network activity",
			Confidence:  0.85,
		}
	}

	// 🟡 Fallback: long idle gap with no useful context
	if gapSeconds >= 2 {
		return Result{
			Cause:       CauseNetwork,
			Character:   "🐌 Network Slug",
			Explanation: "Build was blocked by an external or idle wait",
			Confidence:  0.6,
		}
	}

	return Result{
		Cause:       CauseUnknown,
		Character:   "🤷 Unknown",
		Explanation: "Could not confidently determine the cause",
		Confidence:  0.3,
	}
}

// -----------------------------
// Heuristics
// -----------------------------

func looksLikeNetwork(output string) bool {
	networkKeywords := []string{
		"download",
		"fetch",
		"http",
		"https",
		"registry",
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

func looksLikeCache(output string) bool {
	cacheKeywords := []string{
		"cache",
		"rebuilding",
		"invalidated",
		"cold",
		"node_modules",
		"recompile",
	}

	outputLower := strings.ToLower(output)

	for _, keyword := range cacheKeywords {
		if strings.Contains(outputLower, keyword) {
			return true
		}
	}

	return false
}

func looksLikeDocker(output string) bool {
	dockerKeywords := []string{
		"docker",
		"pull",
		"layer",
		"extracting",
		"image",
		"sha256",
	}

	outputLower := strings.ToLower(output)

	for _, keyword := range dockerKeywords {
		if strings.Contains(outputLower, keyword) {
			return true
		}
	}

	return false
}
