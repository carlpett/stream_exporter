package linemetrics

import (
	"reflect"
	"testing"
	"time"
)

func TestUnmarshalSummaryConfig(t *testing.T) {
	configSlice, err := ReadPatternConfig("testdata/summary.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if len(configSlice) != 1 {
		t.Fatal("Unexpected length of test data")
	}
	config := configSlice[0].SummaryConfig

	expectedMaxAge, _ := time.ParseDuration("15m")
	expectedConfig := SummaryConfig{
		Objectives: map[float64]float64{
			0.1: 0.1,
			0.2: 0.2,
			0.3: 0.3,
			0.4: 0.4,
		},
		MaxAge:     expectedMaxAge,
		AgeBuckets: 3,
		BufCap:     500,
	}

	if !reflect.DeepEqual(expectedConfig, config) {
		t.Fatalf("Mismatch! Expected\n\t%v\nGot\n\t%v\n", expectedConfig, config)
	}
}

func TestUnmarshalHistogramConfig(t *testing.T) {
	configSlice, err := ReadPatternConfig("testdata/histogram.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if len(configSlice) != 1 {
		t.Fatal("Unexpected length of test data")
	}
	config := configSlice[0].HistogramConfig

	expectedConfig := HistogramConfig{
		Buckets: []float64{0.1, 0.2, 0.5},
	}

	if !reflect.DeepEqual(expectedConfig, config) {
		t.Fatalf("Mismatch! Expected\n\t%v\nGot\n\t%v\n", expectedConfig, config)
	}
}
