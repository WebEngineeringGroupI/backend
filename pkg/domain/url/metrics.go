package url

//go:generate mockgen -source=$GOFILE -destination=./mocks/${GOFILE} -package=mocks
type Metrics interface {
	RecordSingleURLMetrics()
	RecordFileURLMetrics()
}
