package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"

	"github.com/carlpett/stream_exporter/input"
	"github.com/carlpett/stream_exporter/linemetrics"
)

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
	Input   input.InputConfig           `mapstructure:"input"`
	Metrics []linemetrics.MetricsConfig `mapstructure:"metrics"`
}

var (
	configFilePath = flag.String("config-file", "stream_exporter.yaml", "path to config file")
)

func main() {
	// Read configuration
	content, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		os.Exit(1)
	}
	rawConfig := make(map[string]interface{})
	err = yaml.Unmarshal(content, &rawConfig)
	if err != nil {
		panic(err)
	}
	var config Config
	err = mapstructure.Decode(rawConfig, &config)
	if err != nil {
		panic(err)
	}

	// Define metrics
	metrics := make([]linemetrics.LineMetric, 0, len(config.Metrics))
	for _, definition := range config.Metrics {
		lineMetric := linemetrics.NewLineMetric(definition.Name, definition.Pattern, definition.Kind)
		metrics = append(metrics, lineMetric)
	}

	prometheus.MustRegister(lineProcessingTime)

	// Setup signal handling
	quitSig := make(chan os.Signal, 1)
	signal.Notify(quitSig, os.Interrupt)

	// Configure input
	inputReader := input.NewInput(config.Input)
	inputChannel := make(chan string)
	go inputReader.StartStream(inputChannel)

	// Main loop
	done := false
	for !done {
		select {
		case line, ok := <-inputChannel:
			if !ok {
				done = true
				break
			}
			for _, m := range metrics {
				t := time.Now()
				m.MatchLine(line)
				lineProcessingTime.WithLabelValues(m.Name()).Observe(time.Since(t).Seconds())
			}
		case <-quitSig:
			fmt.Println("Received quit signal, shutting down...")
			done = true
			break
		}
	}

	metfam, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		fmt.Println(err)
	}
	for _, met := range metfam {
		fmt.Println(met)
	}
}
