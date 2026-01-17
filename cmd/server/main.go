package main

import (
	"log"
	"net/http"
	"time"

	"k8s-poc/internal/config"
	apphttp "k8s-poc/internal/http"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	port := config.GetEnv("PORT", "8080")

	// Registry padrão do Prometheus (pode ser customizado depois)
	reg := prometheus.DefaultRegisterer
	apphttp.RegisterMetrics(reg)

	handlers := apphttp.NewHandlers()
	router := apphttp.NewRouter(handlers, reg)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("server started on :%s (pod=%s env=%s)",
		port,
		config.GetEnv("POD_NAME", "unknown"),
		config.GetEnv("ENVIRONMENT", "unknown"),
	)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
