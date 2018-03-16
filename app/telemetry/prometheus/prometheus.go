package prometheus

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Config holds data structs for prometheus
type Config struct {
	processCount    *prometheus.CounterVec
	processDuration prometheus.Histogram
}

// Init creates and registers prometheus handles
func Init() (*Config, error) {
	c := Config{}

	c.processCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "nixie",
			Subsystem: "telemetry",
			Name:      "process_count",
			Help:      "The number of messages processed.",
		},
		[]string{"state"},
	)

	if err := prometheus.Register(c.processCount); err != nil {
		return nil, err
	}

	c.processDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "nixie",
			Subsystem: "telemetry",
			Name:      "process_duration_seconds",
			Help:      "The process duration in seconds.",
		},
	)

	if err := prometheus.Register(c.processDuration); err != nil {
		return nil, err
	}

	return &c, nil
}

// IncTelemetryProcessCount increases the Telemetry Process counter
func (c *Config) IncTelemetryProcessCount(state string) {
	c.processCount.WithLabelValues(state).Inc()
}

// ObserveProcessDuration logs the elapsed time for processing a task
func (c *Config) ObserveProcessDuration(start int64) {
	elapsed := time.Now().Unix() - start
	c.processDuration.Observe(float64(elapsed))
}

// EOF
