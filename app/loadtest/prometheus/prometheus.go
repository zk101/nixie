package prometheus

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Config holds data structs for prometheus
type Config struct {
	reqCount    *prometheus.CounterVec
	reqDuration *prometheus.HistogramVec
}

// Init creates and registers prometheus handles
func Init() (*Config, error) {
	c := Config{}

	c.reqCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "nixie",
			Subsystem: "loadtest",
			Name:      "request_total",
			Help:      "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"worker", "test", "state"},
	)

	if err := prometheus.Register(c.reqCount); err != nil {
		return nil, err
	}

	c.reqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "nixie",
			Subsystem: "loadtest",
			Name:      "request_duration_seconds",
			Help:      "The HTTP request latencies in seconds.",
		},
		[]string{"test"},
	)

	if err := prometheus.Register(c.reqDuration); err != nil {
		return nil, err
	}

	return &c, nil
}

// IncReqCount does the counter increment for a request
func (c *Config) IncReqCount(worker, test string, status bool) {
	state := "good"
	if status == false {
		state = "bad"
	}
	c.reqCount.WithLabelValues(worker, test, state).Inc()
}

// ObserveReqDuration determines the elapsed time for a request
func (c *Config) ObserveReqDuration(start *time.Time, test string) {
	elapsed := float64(time.Since(*start)) / float64(time.Second)
	c.reqDuration.WithLabelValues(test).Observe(elapsed)
}

// EOF
