package linemetrics

import (
	"errors"
)

type LineMetric interface {
	MatchLine(s string)
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
