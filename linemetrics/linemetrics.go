package linemetrics

import (
	"errors"
	"io/ioutil"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
)

type LineMetric interface {
	MatchLine(s string)
	Name() string
}

type BaseLineMetric struct {
	name    string
	pattern regexp.Regexp
	labels  []string
}

func (m BaseLineMetric) Name() string {
	return m.name
}

type config struct {
	Metrics []MetricsConfig
}

func ReadPatternConfig(path string) ([]MetricsConfig, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return config.Metrics, nil
}

func NewLineMetric(name string, rawPattern string, kind metricKind) (LineMetric, prometheus.Collector) {
	pattern := regexp.MustCompile(rawPattern)
	labels := pattern.SubexpNames()[1:] // First element is entire expression

	var lineMetric LineMetric
	base := BaseLineMetric{
		name:    name,
		pattern: *pattern,
		labels:  labels,
	}
	var collector prometheus.Collector
	switch kind {
	case counter:
		lineMetric, collector = NewCounterLineMetric(base)
	case gauge:
		lineMetric, collector = NewGaugeLineMetric(base)
	case histogram:
		lineMetric, collector = NewHistogramLineMetric(base)
	case summary:
		lineMetric, collector = NewSummaryLineMetric(base)
	}

	return lineMetric, collector
}

func getValueCaptureIndex(labels []string) (int, error) {
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
		return valueIdx, errors.New("No named capture group for 'value'")
	}

	return valueIdx, nil
}
