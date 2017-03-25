package linemetrics

import (
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
)

type CounterVecLineMetric struct {
	pattern regexp.Regexp
	labels  []string
	metric  prometheus.CounterVec
}

func (counter CounterVecLineMetric) MatchLine(s string) {
	matches := counter.pattern.FindStringSubmatch(s)
	if len(matches) > 0 {
		captures := matches[1:]
		counter.metric.WithLabelValues(captures...).Inc()
	}
}

type CounterLineMetric struct {
	pattern regexp.Regexp
	metric  prometheus.Counter
}

func (counter CounterLineMetric) MatchLine(s string) {
	matches := counter.pattern.MatchString(s)
	if matches {
		counter.metric.Inc()
	}
}

func NewCounterLineMetric(name string, labels []string, pattern *regexp.Regexp) LineMetric {
	opts := prometheus.CounterOpts{
		Name: name,
		Help: name,
	}
	var lineMetric LineMetric
	if len(labels) > 0 {
		metric := prometheus.NewCounterVec(opts, labels)
		lineMetric = CounterVecLineMetric{
			pattern: *pattern,
			labels:  labels,
			metric:  *metric,
		}
		prometheus.Register(metric)

	} else {
		metric := prometheus.NewCounter(opts)
		lineMetric = CounterLineMetric{
			pattern: *pattern,
			metric:  metric,
		}
		prometheus.Register(metric)
	}

	return lineMetric
}
