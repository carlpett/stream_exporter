package linemetrics

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"

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

func NewLineMetric(name string, rawPattern string, kind metricKind) LineMetric {
	pattern := regexp.MustCompile(rawPattern)
	labels := pattern.SubexpNames()[1:] // First element is entire expression

	var lineMetric LineMetric
	base := BaseLineMetric{
		name:    name,
		pattern: *pattern,
		labels:  labels,
	}
	switch kind {
	case counter:
		lineMetric = NewCounterLineMetric(base)
	case gauge:
		lineMetric = NewGaugeLineMetric(base)
	case histogram:
		lineMetric = NewHistogramLineMetric(base)
	case summary:
		lineMetric = NewSummaryLineMetric(base)
	}

	return lineMetric
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
