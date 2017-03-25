package linemetrics

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type GaugeVecLineMetric struct {
	pattern  regexp.Regexp
	valueIdx int
	labels   []string
	metric   prometheus.GaugeVec
}

func (gauge GaugeVecLineMetric) MatchLine(s string) {
	matches := gauge.pattern.FindStringSubmatch(s)
	if len(matches) > 0 {
		captures := matches[1:]
		valueStr := captures[gauge.valueIdx]
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			fmt.Printf("Unable to convert %s to float\n", valueStr)
			return
		}
		gauge.metric.WithLabelValues(captures...).Set(value)
	}
}

type GaugeLineMetric struct {
	pattern regexp.Regexp
	metric  prometheus.Gauge
}

func (gauge GaugeLineMetric) MatchLine(s string) {
	matches := gauge.pattern.FindStringSubmatch(s)
	if len(matches) > 0 {
		captures := matches[1:]
		valueStr := captures[0] // There are no other labels
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			fmt.Printf("Unable to convert %s to float\n", valueStr)
			return
		}

		gauge.metric.Set(value)
	}
}

func NewGaugeLineMetric(name string, labels []string, pattern *regexp.Regexp) LineMetric {
	valueIdx, err := getValueCaptureIndex(labels)
	if err != nil {
		panic(fmt.Sprintf("Error initializing gauge %s: %s", name, err))
	}
	labels = append(labels[:valueIdx], labels[valueIdx+1:]...)

	opts := prometheus.GaugeOpts{
		Name: name,
		Help: name,
	}
	var lineMetric LineMetric
	if len(labels) > 0 {
		metric := prometheus.NewGaugeVec(opts, labels)
		lineMetric = GaugeVecLineMetric{
			pattern:  *pattern,
			valueIdx: valueIdx,
			labels:   labels,
			metric:   *metric,
		}
		prometheus.Register(metric)

	} else {
		metric := prometheus.NewGauge(opts)
		lineMetric = GaugeLineMetric{
			pattern: *pattern,
			metric:  metric,
		}
		prometheus.Register(metric)
	}

	return lineMetric
}
