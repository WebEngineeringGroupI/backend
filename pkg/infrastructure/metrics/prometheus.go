package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Definir urlShorted
// Crear constructor y asignar valor a urlShorted

/*
= promauto.NewCounter(prometheus.CounterOpts{
		Name: "urlshortener_shorted_url_total",
		Help: "The total number of shorted urls",
		})
 */

type CustomMetrics struct {
	urlShorted prometheus.Counter
}

func NewCustomMetrics() CustomMetrics {
	var urlShorted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "urlshortener_shorted_url_total",
		Help: "The total number of shorted urls",
	})

	return CustomMetrics{urlShorted: urlShorted}
}

// RecordUrlShorted Call this function every time a URL is shortened
func (m *CustomMetrics) RecordUrlShorted() {
	go func() {
		m.urlShorted.Inc()
	}()
}