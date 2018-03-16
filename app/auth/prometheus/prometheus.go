package prometheus

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Config holds data structs for prometheus
type Config struct {
	reqCount       *prometheus.CounterVec
	reqDuration    *prometheus.HistogramVec
	ldapPoolQueue  prometheus.Gauge
	ldapPoolWorker prometheus.Gauge
	cbPoolQueue    prometheus.Gauge
	cbPoolWorker   prometheus.Gauge
}

// Init creates and registers prometheus handles
func Init() (*Config, error) {
	c := Config{}

	c.reqCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "nixie",
			Subsystem: "auth",
			Name:      "request_total",
			Help:      "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"code", "method", "path"},
	)

	if err := prometheus.Register(c.reqCount); err != nil {
		return nil, err
	}

	c.reqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "nixie",
			Subsystem: "auth",
			Name:      "request_duration_seconds",
			Help:      "The HTTP request latencies in seconds.",
		},
		[]string{"path"},
	)

	if err := prometheus.Register(c.reqDuration); err != nil {
		return nil, err
	}

	c.ldapPoolQueue = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "nixie",
			Subsystem: "auth",
			Name:      "ldappool_queue_count",
			Help:      "The number of current queued tasks.",
		},
	)

	if err := prometheus.Register(c.ldapPoolQueue); err != nil {
		return nil, err
	}

	c.ldapPoolWorker = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "nixie",
			Subsystem: "auth",
			Name:      "ldappool_worker_count",
			Help:      "The number of workers running.",
		},
	)

	if err := prometheus.Register(c.ldapPoolWorker); err != nil {
		return nil, err
	}

	c.cbPoolQueue = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "nixie",
			Subsystem: "auth",
			Name:      "couchbasepool_queue_count",
			Help:      "The number of current queued tasks.",
		},
	)

	if err := prometheus.Register(c.cbPoolQueue); err != nil {
		return nil, err
	}

	c.cbPoolWorker = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "nixie",
			Subsystem: "auth",
			Name:      "couchbasepool_worker_count",
			Help:      "The number of workers running.",
		},
	)

	if err := prometheus.Register(c.cbPoolWorker); err != nil {
		return nil, err
	}

	return &c, nil
}

// IncReqCount does the counter increment for a request
func (c *Config) IncReqCount(code int, method, path string) {
	c.reqCount.WithLabelValues(strconv.Itoa(code), method, path).Inc()
}

// ObserveReqDuration determines the elapsed time for a request
func (c *Config) ObserveReqDuration(start *time.Time, path string) {
	elapsed := float64(time.Since(*start)) / float64(time.Second)
	c.reqDuration.WithLabelValues(path).Observe(elapsed)
}

// SetLDAPQueueCount sets the queue gauge to count
func (c *Config) SetLDAPQueueCount(count int) {
	c.ldapPoolQueue.Set(float64(count))
}

// SetLDAPWorkerCount sets the worker gauge to count
func (c *Config) SetLDAPWorkerCount(count int) {
	c.ldapPoolWorker.Set(float64(count))
}

// SetCBQueueCount sets the queue gauge to count
func (c *Config) SetCBQueueCount(count int) {
	c.cbPoolQueue.Set(float64(count))
}

// SetCBWorkerCount sets the worker gauge to count
func (c *Config) SetCBWorkerCount(count int) {
	c.cbPoolWorker.Set(float64(count))
}

// EOF
