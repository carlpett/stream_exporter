package linemetrics

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type SummaryLineMetric struct {
	pattern  regexp.Regexp
	valueIdx int
	metric   prometheus.Summary
}

func (summary SummaryLineMetric) MatchLine(s string) {
	matches := summary.pattern.FindStringSubmatch(s)
	if len(matches) > 0 {
		captures := matches[1:]
		valueStr := captures[summary.valueIdx]
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			fmt.Printf("Unable to convert %s to float\n", valueStr)
			return
		}
		summary.metric.Observe(value)
	}
}

type SummaryVecLineMetric struct {
	pattern  regexp.Regexp
	labels   []string
	valueIdx int
	metric   prometheus.SummaryVec
}

func (summary SummaryVecLineMetric) MatchLine(s string) {
	matches := summary.pattern.FindStringSubmatch(s)
	if len(matches) > 0 {
		captures := matches[1:]
		valueStr := captures[summary.valueIdx]
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			fmt.Printf("Unable to convert %s to float\n", valueStr)
			return
		}
		capturedLabels := append(captures[0:summary.valueIdx], captures[summary.valueIdx+1:]...)
		summary.metric.WithLabelValues(capturedLabels...).Observe(value)
	}
}

func NewSummaryLineMetric(name string, labels []string, pattern *regexp.Regexp) LineMetric {
	valueIdx, err := getValueCaptureIndex(labels)
	if err != nil {
		panic(fmt.Sprintf("Error initializing summary %s: %s", name, err))
	}
	labels = append(labels[:valueIdx], labels[valueIdx+1:]...)

	opts := prometheus.SummaryOpts{
		Name: name,
		Help: name,
	}
	var lineMetric LineMetric
	if len(labels) > 0 {
		metric := prometheus.NewSummaryVec(opts, labels)
		lineMetric = SummaryVecLineMetric{
			pattern:  *pattern,
			labels:   labels,
			valueIdx: valueIdx,
			metric:   *metric,
		}
		prometheus.Register(metric)
	} else {
		metric := prometheus.NewSummary(opts)
		lineMetric = SummaryLineMetric{
			pattern:  *pattern,
			valueIdx: valueIdx,
			metric:   metric,
		}
		prometheus.Register(metric)
	}

	return lineMetric
}
