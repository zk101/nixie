package prometheus

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Direction constants for defining well, direction
const (
	DirectionRX = "rx"
	DirectionTX = "tx"
)

// Config holds data structs for prometheus
type Config struct {
	msgCount     *prometheus.CounterVec
	msgDuration  *prometheus.HistogramVec
	msgSize      *prometheus.HistogramVec
	connCount    prometheus.Gauge
	connTotal    prometheus.Counter
	connDuration prometheus.Histogram
	cbPoolQueue  prometheus.Gauge
	cbPoolWorker prometheus.Gauge
	goPoolQueue  prometheus.Gauge
	goPoolWorker prometheus.Gauge

	// async pool
	asynctxQueue    prometheus.Gauge
	asynctxWorker   prometheus.Gauge
	asynctxDelivery *prometheus.CounterVec
	asyncrxWorker   prometheus.Gauge
	asyncrxDelivery *prometheus.CounterVec
}

// Init creates and registers prometheus handles
func Init() (*Config, error) {
	c := Config{}

	c.msgCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "nixie",
			Subsystem: "ws",
			Name:      "message_total",
			Help:      "The number of messages processed",
		},
		[]string{"msg_type", "direction", "security"},
	)

	if err := prometheus.Register(c.msgCount); err != nil {
		return nil, err
	}

	c.msgDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "nixie",
			Subsystem: "ws",
			Name:      "message_process_seconds",
			Help:      "The message processing latencies in seconds.",
		},
		[]string{"msg_type", "direction", "security"},
	)

	if err := prometheus.Register(c.msgDuration); err != nil {
		return nil, err
	}

	c.msgSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "nixie",
			Subsystem: "ws",
			Name:      "message_size_bytes",
			Help:      "The message size in bytes.",
		},
		[]string{"msg_type", "direction", "security"},
	)

	if err := prometheus.Register(c.msgSize); err != nil {
		return nil, err
	}

	c.connCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "nixie",
			Subsystem: "ws",
			Name:      "connection_current_count",
			Help:      "The number of current connection.",
		},
	)

	if err := prometheus.Register(c.connCount); err != nil {
		return nil, err
	}

	c.connTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "nixie",
			Subsystem: "ws",
			Name:      "connection_total",
			Help:      "The total number of connections.",
		},
	)

	if err := prometheus.Register(c.connTotal); err != nil {
		return nil, err
	}

	c.connDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "nixie",
			Subsystem: "ws",
			Name:      "connection_duration_seconds",
			Help:      "The connection duration in seconds.",
		},
	)

	if err := prometheus.Register(c.connDuration); err != nil {
		return nil, err
	}

	c.cbPoolQueue = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "nixie",
			Subsystem: "ws",
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
			Subsystem: "ws",
			Name:      "couchbasepool_worker_count",
			Help:      "The number of workers running.",
		},
	)

	if err := prometheus.Register(c.cbPoolWorker); err != nil {
		return nil, err
	}

	c.goPoolQueue = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "nixie",
			Subsystem: "ws",
			Name:      "gopool_queue_count",
			Help:      "The number of current queued tasks.",
		},
	)

	if err := prometheus.Register(c.goPoolQueue); err != nil {
		return nil, err
	}

	c.goPoolWorker = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "nixie",
			Subsystem: "ws",
			Name:      "gopool_worker_count",
			Help:      "The number of workers running.",
		},
	)

	if err := prometheus.Register(c.goPoolWorker); err != nil {
		return nil, err
	}

	c.asynctxQueue = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "nixie",
			Subsystem: "ws",
			Name:      "asynctx_queue_count",
			Help:      "The number of queued tasks.",
		},
	)

	if err := prometheus.Register(c.asynctxQueue); err != nil {
		return nil, err
	}

	c.asynctxWorker = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "nixie",
			Subsystem: "ws",
			Name:      "asynctx_worker_count",
			Help:      "The number of workers running.",
		},
	)

	if err := prometheus.Register(c.asynctxWorker); err != nil {
		return nil, err
	}

	c.asynctxDelivery = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "nixie",
			Subsystem: "ws",
			Name:      "asynctx_delivery_count",
			Help:      "The number of messages delivered",
		},
		[]string{"queue", "state"},
	)

	if err := prometheus.Register(c.asynctxDelivery); err != nil {
		return nil, err
	}

	c.asyncrxWorker = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "nixie",
			Subsystem: "ws",
			Name:      "asyncrx_worker_count",
			Help:      "The number of workers running.",
		},
	)

	if err := prometheus.Register(c.asyncrxWorker); err != nil {
		return nil, err
	}

	c.asyncrxDelivery = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "nixie",
			Subsystem: "ws",
			Name:      "asyncrx_delivery_count",
			Help:      "The number of messages delivered",
		},
		[]string{"queue", "state"},
	)

	if err := prometheus.Register(c.asyncrxDelivery); err != nil {
		return nil, err
	}

	return &c, nil
}

// IncReqCount does the counter increment for a request
func (c *Config) IncReqCount(msgType, direction, security string) {
	c.msgCount.WithLabelValues(msgType, direction, security).Inc()
}

// ObserveReqDuration logs the elapsed time for a request
func (c *Config) ObserveReqDuration(start time.Time, msgType, direction, security string) {
	elapsed := float64(time.Since(start)) / float64(time.Second)
	c.msgDuration.WithLabelValues(msgType, direction, security).Observe(elapsed)
}

// ObserveReqSize logs the size for a request
func (c *Config) ObserveReqSize(size int64, msgType, direction, security string) {
	c.msgSize.WithLabelValues(msgType, direction, security).Observe(float64(size))
}

// SetConnCount moves the connection count gauge up and down
func (c *Config) SetConnCount(count int) {
	c.connCount.Set(float64(count))
}

// IncConnTotal does the counter increment for a connection
func (c *Config) IncConnTotal() {
	c.connTotal.Inc()
}

// ObserveConnDuration logs the elapsed time for a connection
func (c *Config) ObserveConnDuration(start int64) {
	elapsed := time.Now().Unix() - start
	c.connDuration.Observe(float64(elapsed))
}

// SetCBQueueCount sets the queue gauage to count
func (c *Config) SetCBQueueCount(count int) {
	c.cbPoolQueue.Set(float64(count))
}

// SetCBWorkerCount sets the worker gauage to count
func (c *Config) SetCBWorkerCount(count int) {
	c.cbPoolWorker.Set(float64(count))
}

// SetGoPoolQueueCount sets the queue gauage to count
func (c *Config) SetGoPoolQueueCount(count int) {
	c.goPoolQueue.Set(float64(count))
}

// SetGoPoolWorkerCount sets the worker gauage to count
func (c *Config) SetGoPoolWorkerCount(count int) {
	c.goPoolWorker.Set(float64(count))
}

// SetAsynctxQueueCount sets the queue gauage to count
func (c *Config) SetAsynctxQueueCount(count int) {
	c.asynctxQueue.Set(float64(count))
}

// SetAsynctxWorkerCount sets the worker gauage to count
func (c *Config) SetAsynctxWorkerCount(count int) {
	c.asynctxWorker.Set(float64(count))
}

// IncAsynctxDeliveryCount increases the Async TX delivery counter
func (c *Config) IncAsynctxDeliveryCount(queue, state string) {
	c.asynctxDelivery.WithLabelValues(queue, state).Inc()
}

// SetAsyncrxWorkerCount sets the worker gauage to count
func (c *Config) SetAsyncrxWorkerCount(count int) {
	c.asyncrxWorker.Set(float64(count))
}

// IncAsyncrxDeliveryCount increases the Async RX delivery counter
func (c *Config) IncAsyncrxDeliveryCount(queue, state string) {
	c.asyncrxDelivery.WithLabelValues(queue, state).Inc()
}

// EOF
