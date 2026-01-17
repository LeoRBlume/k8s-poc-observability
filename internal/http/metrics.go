package http

import (
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	healthRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "health_requests_total",
			Help: "Total de requisicoes processadas pelo endpoint /health.",
		},
		[]string{"pod", "remote_ip"},
	)

	healthRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "health_request_duration_seconds",
			Help:    "Duracao do handler /health em segundos.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"pod"},
	)
)

func RegisterMetrics(reg prometheus.Registerer) {
	reg.MustRegister(healthRequestsTotal, healthRequestDuration)
}

func HealthMetricsMiddleware(podName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Instrumenta somente /health
		if c.FullPath() != "/health" {
			c.Next()
			return
		}

		start := time.Now()
		remoteIP := extractRemoteIP(c.Request.RemoteAddr)

		c.Next()

		healthRequestsTotal.WithLabelValues(podName, remoteIP).Inc()
		healthRequestDuration.WithLabelValues(podName).Observe(time.Since(start).Seconds())
	}
}

func extractRemoteIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return host
}
