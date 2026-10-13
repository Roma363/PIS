package src

// ProbeResult містить результат опитування ендпоінту
type ProbeResult struct {
	URL        string
	StatusCode int
	LatencyMs  int64
	Status     string
}

// EvaluateStatus обчислює доступність згідно з лімітом затримки та кодом
func EvaluateStatus(statusCode int, latencyMs int64, timeoutMs int64) string {
	if latencyMs > timeoutMs {
		return "DOWN"
	}
	if statusCode >= 200 && statusCode < 400 {
		return "UP"
	}
	return "DOWN"
}
