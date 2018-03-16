package prometheus

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Config holds data structs for prometheus
type Config struct {
	processCount    *prometheus.CounterVec
	processDuration *prometheus.HistogramVec
	ldapPoolQueue   prometheus.Gauge
	ldapPoolWorker  prometheus.Gauge
	cbPoolQueue     prometheus.Gauge
	cbPoolWorker    prometheus.Gauge
}

// Init creates and registers prometheus handles
func Init() (*Config, error) {
	c := Config{}

	c.processCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "nixie",
			Subsystem: "chat",
			Name:      "process_count",
			Help:      "The number of messages processed.",
		},
		[]string{"state", "category"},
	)

	if err := prometheus.Register(c.processCount); err != nil {
		return nil, err
	}

	c.processDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "nixie",
			Subsystem: "chat",
			Name:      "process_duration_seconds",
			Help:      "The process duration in seconds.",
		},
		[]string{"state", "category"},
	)

	if err := prometheus.Register(c.processDuration); err != nil {
		return nil, err
	}

	c.ldapPoolQueue = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "nixie",
			Subsystem: "chat",
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
			Subsystem: "chat",
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
			Subsystem: "chat",
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
			Subsystem: "chat",
			Name:      "couchbasepool_worker_count",
			Help:      "The number of workers running.",
		},
	)

	if err := prometheus.Register(c.cbPoolWorker); err != nil {
		return nil, err
	}

	return &c, nil
}

// IncChatProcessCount increases the Chat Process counter
func (c *Config) IncChatProcessCount(state, category string) {
	c.processCount.WithLabelValues(state, category).Inc()
}

// ObserveProcessDuration logs the elapsed time for processing a task
func (c *Config) ObserveProcessDuration(start int64, state, category string) {
	elapsed := time.Now().Unix() - start
	c.processDuration.WithLabelValues(state, category).Observe(float64(elapsed))
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
