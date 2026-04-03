package metrics

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	apiRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "uptime_api_requests_total",
			Help: "Total number of API requests",
		},
		[]string{"method", "path", "status"},
	)

	checksTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "uptime_checks_total",
			Help: "Total number of uptime checks executed",
		},
		[]string{"result", "status_code"},
	)

	checkLatencyMs = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "uptime_check_latency_ms",
			Help:    "Latency of uptime checks in milliseconds",
			Buckets: []float64{50, 100, 250, 500, 1000, 2000, 5000, 10000},
		},
	)

	enabledTargets = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "uptime_enabled_targets",
			Help: "Current number of enabled targets",
		},
	)
)

func init() {
	prometheus.MustRegister(apiRequestsTotal, checksTotal, checkLatencyMs, enabledTargets)
}

func Handler() http.Handler {
	return promhttp.Handler()
}

func RecordAPIRequest(method, path string, status int) {
	apiRequestsTotal.WithLabelValues(method, path, strconv.Itoa(status)).Inc()
}

func RecordCheck(success bool, statusCode int, latencyMs float64) {
	result := "failure"
	if success {
		result = "success"
	}
	checksTotal.WithLabelValues(result, strconv.Itoa(statusCode)).Inc()
	checkLatencyMs.Observe(latencyMs)
}

func SetEnabledTargets(n int) {
	enabledTargets.Set(float64(n))
}
