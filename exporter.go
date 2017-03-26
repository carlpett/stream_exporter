package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"

	"github.com/carlpett/stream_exporter/linemetrics"
)

func rmfifo() {
	os.Remove("/tmp/myfifo")
}

var (
	lineProcessingTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "stream_exporter",
			Subsystem: "line_processing",
			Name:      "duration_seconds",
			Help:      "Observed duration, in seconds, of processing a single line per registered metric",
			Buckets:   prometheus.ExponentialBuckets(time.Microsecond.Seconds(), 3.981072, 5),
			// This results in 5 buckets from 1 us to approx 1 ms (3.98...^5 ~= 1000)
		},
		[]string{"metric"},
	)
)

type Config struct {
	Metrics []linemetrics.MetricsConfig
}

func main() {
	// "Read config"
	content, err := ioutil.ReadFile("test-config.yaml")
	if err != nil {
		panic(err)
	}
	config := Config{}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		panic(err)
	}

	// "Define metrics"
	metrics := make([]linemetrics.LineMetric, 0, len(config.Metrics))
	for _, definition := range config.Metrics {
		lineMetric := linemetrics.NewLineMetric(definition.Name, definition.Pattern, definition.Kind)
		metrics = append(metrics, lineMetric)
	}

	prometheus.MustRegister(lineProcessingTime)

	// "Config input"
	err = syscall.Mkfifo("/tmp/myfifo", 0666)
	if err != nil {
		panic(err)
	}
	defer rmfifo()

	pipe, err := os.OpenFile("/tmp/myfifo", os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(pipe)
	scanner := bufio.NewScanner(reader)

	// "Main loop"
	for scanner.Scan() {
		line := scanner.Text()
		for _, m := range metrics {
			t := time.Now()
			m.MatchLine(line)
			lineProcessingTime.WithLabelValues(m.Name()).Observe(time.Since(t).Seconds())
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	metfam, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		fmt.Println(err)
	}
	for _, met := range metfam {
		fmt.Println(met)
	}
}
