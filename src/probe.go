package src

// ProbeResult містить результат опитування ендпоінту
type ProbeResult struct {
	URL        string
	StatusCode int
	LatencyMs  int64
	Status     string
}

func EvaluateStatus(statusCode int, latencyMs int64, timeoutMs int64) string {
	if statusCode == 200 {
		return "UP"
	}

	return "DOWN"
}
