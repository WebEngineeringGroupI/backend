package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type PrometheusMetrics struct {
	urlsProcessed         prometheus.Counter
	singleUrlsProcessed   prometheus.Counter
	multipleUrlsProcessed prometheus.Counter
	fileUrlsProcessed     prometheus.Counter
}

func NewPrometheusMetrics() *PrometheusMetrics {
	var urlsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "urlshortener_urls_processed_total",
		Help: "The total number of shorted urls",
	})

	var singleUrlsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "urlshortener_single_urls_processed_total",
		Help: "The total number of single shorted urls",
	})

	var multipleUrlsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "urlshortener_multiple_urls_processed_total",
		Help: "The total number of multiple shorted urls",
	})

	var fileUrlsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "urlshortener_file_urls_processed_total",
		Help: "The total number of file shorted urls",
	})

	return &PrometheusMetrics{
		urlsProcessed:         urlsProcessed,
		singleUrlsProcessed:   singleUrlsProcessed,
		multipleUrlsProcessed: multipleUrlsProcessed,
		fileUrlsProcessed:     fileUrlsProcessed,
	}
}

func (r *PrometheusMetrics) RecordUrlsProcessed() {
	r.urlsProcessed.Inc()
}

func (r *PrometheusMetrics) RecordSingleURLMetrics() {
	r.singleUrlsProcessed.Inc()
}

func (r *PrometheusMetrics) RecordMultipleURLMetrics() {
	r.multipleUrlsProcessed.Inc()
}

func (r *PrometheusMetrics) RecordFileURLMetrics() {
	r.fileUrlsProcessed.Inc()
}
