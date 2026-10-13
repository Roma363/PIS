package src

// ProbeResult містить результат опитування ендпоінту
type ProbeResult struct {
	URL        string
	StatusCode int
	LatencyMs  int64
	Status     string
}
