package linemetrics

type LineMetric interface {
	MatchLine(s string)
}
