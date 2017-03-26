package linemetrics

type metricKind string

const (
	counter   metricKind = "counter"
	gauge     metricKind = "gauge"
	histogram metricKind = "histogram"
	summary   metricKind = "summary"
)

type MetricsConfig struct {
	Name    string
	Kind    metricKind
	Pattern string
}
