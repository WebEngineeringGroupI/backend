package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {
	factory := newFactory()

	// This two instructions expose application metric on http://localhost:2112/metrics
	// We can expose our own metrics, more info here https://prometheus.io/docs/guides/go-application/#adding-your-own-metrics
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)

	log.Fatal(http.ListenAndServe(":8080", factory.NewHTTPRouter()))
}
