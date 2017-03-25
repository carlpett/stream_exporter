package linemetrics

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type HistogramLineMetric struct {
	pattern  regexp.Regexp
	valueIdx int
	metric   prometheus.Histogram
}

func (histogram HistogramLineMetric) MatchLine(s string) {
	matches := histogram.pattern.FindStringSubmatch(s)
	if len(matches) > 0 {
		captures := matches[1:]
		valueStr := captures[histogram.valueIdx]
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			fmt.Printf("Unable to convert %s to float\n", valueStr)
			return
		}
		histogram.metric.Observe(value)
	}
}

type HistogramVecLineMetric struct {
	pattern  regexp.Regexp
	labels   []string
	valueIdx int
	metric   prometheus.HistogramVec
}

func (histogram HistogramVecLineMetric) MatchLine(s string) {
	matches := histogram.pattern.FindStringSubmatch(s)
	if len(matches) > 0 {
		captures := matches[1:]
		valueStr := captures[histogram.valueIdx]
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			fmt.Printf("Unable to convert %s to float\n", valueStr)
			return
		}
		capturedLabels := append(captures[0:histogram.valueIdx], captures[histogram.valueIdx+1:]...)
		histogram.metric.WithLabelValues(capturedLabels...).Observe(value)
	}
}

func NewHistogramLineMetric(name string, labels []string, pattern *regexp.Regexp) LineMetric {
	valueIdx, err := getValueCaptureIndex(labels)
	if err != nil {
		panic(fmt.Sprintf("Error initializing histogram %s: %s", name, err))
	}
	labels = append(labels[:valueIdx], labels[valueIdx+1:]...)

	opts := prometheus.HistogramOpts{
		Name: name,
		Help: name,
	}
	var lineMetric LineMetric
	if len(labels) > 0 {
		metric := prometheus.NewHistogramVec(opts, labels)
		lineMetric = HistogramVecLineMetric{
			pattern:  *pattern,
			labels:   labels,
			valueIdx: valueIdx,
			metric:   *metric,
		}
		prometheus.Register(metric)
	} else {
		metric := prometheus.NewHistogram(opts)
		lineMetric = HistogramLineMetric{
			pattern:  *pattern,
			valueIdx: valueIdx,
			metric:   metric,
		}
		prometheus.Register(metric)
	}

	return lineMetric
}
