package http

import (
	"k8s-poc/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(h *Handlers, reg prometheus.Registerer) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())

	podName := config.GetEnv("POD_NAME", "unknown")
	r.Use(HealthMetricsMiddleware(podName))

	r.GET("/health", h.Health)
	r.GET("/whoami", h.WhoAmI)

	// Expor métricas Prometheus
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}
