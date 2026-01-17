package http

import (
	"net/http"
	"os"
	"time"

	"k8s-poc/internal/config"

	"github.com/gin-gonic/gin"
)

type WhoAmIResponse struct {
	Environment string `json:"environment"`

	PodName   string `json:"podName"`
	PodIP     string `json:"podIP"`
	NodeName  string `json:"nodeName"`
	Namespace string `json:"namespace"`

	Hostname string `json:"hostname"`
	TimeUTC  string `json:"timeUtc"`
}

type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"pod":    config.GetEnv("POD_NAME", "unknown"),
	})
}

func (h *Handlers) WhoAmI(c *gin.Context) {
	hostname, _ := os.Hostname()

	resp := WhoAmIResponse{
		Environment: config.GetEnv("ENVIRONMENT", "unknown"),

		PodName:   config.GetEnv("POD_NAME", "unknown"),
		PodIP:     config.GetEnv("POD_IP", "unknown"),
		NodeName:  config.GetEnv("NODE_NAME", "unknown"),
		Namespace: config.GetEnv("POD_NAMESPACE", "default"),

		Hostname: hostname,
		TimeUTC:  time.Now().UTC().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, resp)
}
