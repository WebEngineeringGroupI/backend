package url

type Metrics interface {
	RecordSingleURLMetrics()
	RecordMultipleURLMetrics()
	RecordFileURLMetrics()
	RecordUrlsProcessed()
}
