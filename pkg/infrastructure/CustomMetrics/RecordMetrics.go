package CustomMetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type RecordMetrics struct {
	urlsProcessed         prometheus.Counter
	singleUrlsProcessed   prometheus.Counter
	multipleUrlsProcessed prometheus.Counter
	fileUrlsProcessed     prometheus.Counter
}

func NewCustomMetrics() RecordMetrics {
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

	return RecordMetrics{
		urlsProcessed:         urlsProcessed,
		singleUrlsProcessed:   singleUrlsProcessed,
		multipleUrlsProcessed: multipleUrlsProcessed,
		fileUrlsProcessed:     fileUrlsProcessed,
	}
}

func (r *RecordMetrics) RecordUrlsProcessed() {
	go func() {
		r.urlsProcessed.Inc()
	}()
}

func (r *RecordMetrics) RecordSingleURLMetrics() {
	go func() {
		r.singleUrlsProcessed.Inc()
	}()
}

func (r *RecordMetrics) RecordMultipleURLMetrics() {
	go func() {
		r.multipleUrlsProcessed.Inc()
	}()
}

func (r *RecordMetrics) RecordFileURLMetrics() {
	go func() {
		r.fileUrlsProcessed.Inc()
	}()
}
