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
		valueStr := matches[histogram.valueIdx+1]
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
		valueStr := matches[histogram.valueIdx+1]
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			fmt.Printf("Unable to convert %s to float\n", valueStr)
			fmt.Printf("Matches: %v, Idx: %d\n", matches, histogram.valueIdx)
			return
		}
		capturedLabels := append(matches[1:histogram.valueIdx+1], matches[histogram.valueIdx+2:]...)
		fmt.Println(capturedLabels)
		histogram.metric.WithLabelValues(capturedLabels...).Observe(value)
	}
}

func NewHistogramLineMetric(name string, labels []string, pattern *regexp.Regexp) LineMetric {
	foundValue := false
	valueIdx := 0
	for idx, l := range labels {
		if l == "value" {
			foundValue = true
			valueIdx = idx
			break
		}
	}
	if !foundValue {
		panic("No capture group for value in histogram")
	}
	labels = append(labels[0:valueIdx], labels[valueIdx+1:]...)

	opts := prometheus.HistogramOpts{
		Name: name,
		Help: name,
	}
	var lineMetric LineMetric
	if len(labels) > 1 {
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
