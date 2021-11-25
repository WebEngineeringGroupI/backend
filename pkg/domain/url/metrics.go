package url

type Metrics interface {
	RecordSingleURLMetrics()
	RecordFileURLMetrics()
	RecordUrlsProcessed()
}
