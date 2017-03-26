package linemetrics

import (
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type GaugeVecLineMetric struct {
	BaseLineMetric
	valueIdx int
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
	BaseLineMetric
	metric prometheus.Gauge
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

func NewGaugeLineMetric(base BaseLineMetric) LineMetric {
	valueIdx, err := getValueCaptureIndex(base.labels)
	if err != nil {
		panic(fmt.Sprintf("Error initializing gauge %s: %s", base.name, err))
	}
	base.labels = append(base.labels[:valueIdx], base.labels[valueIdx+1:]...)

	opts := prometheus.GaugeOpts{
		Name: base.name,
		Help: base.name,
	}
	var lineMetric LineMetric
	if len(base.labels) > 0 {
		metric := prometheus.NewGaugeVec(opts, base.labels)
		lineMetric = GaugeVecLineMetric{
			BaseLineMetric: base,
			valueIdx:       valueIdx,
			metric:         *metric,
		}
		prometheus.Register(metric)

	} else {
		metric := prometheus.NewGauge(opts)
		lineMetric = GaugeLineMetric{
			BaseLineMetric: base,
			metric:         metric,
		}
		prometheus.Register(metric)
	}

	return lineMetric
}
