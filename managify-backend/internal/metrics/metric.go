package metrics

import "time"

type EndpointMetrics struct {
	Count     int64
	TotalTime time.Duration
	MinTime   time.Duration
	MaxTime   time.Duration
}

var Metrics = make(map[string]*EndpointMetrics)
