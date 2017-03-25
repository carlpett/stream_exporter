package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"

	"github.com/carlpett/stream_exporter/linemetrics"
)

func rmfifo() {
	os.Remove("/tmp/myfifo")
}

type metricKind string

const (
	counter   metricKind = "counter"
	gauge     metricKind = "gauge"
	histogram metricKind = "histogram"
	summary   metricKind = "summary"
)

type Config struct {
	Metrics []struct {
		Name    string
		Kind    metricKind
		Pattern string
	}
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

	// "Define metrics"
	metrics := make([]linemetrics.LineMetric, 0, len(config.Metrics))
	for _, definition := range config.Metrics {
		pattern := regexp.MustCompile(definition.Pattern)
		labels := pattern.SubexpNames()[1:] // First element is entire expression

		switch definition.Kind {
		case counter:
			lineMetric := linemetrics.NewCounterLineMetric(definition.Name, labels, pattern)
			metrics = append(metrics, lineMetric)
		case gauge:
		case histogram:
			lineMetric := linemetrics.NewHistogramLineMetric(definition.Name, labels, pattern)
			metrics = append(metrics, lineMetric)
		case summary:
		}
	}

	// "Main loop"
	for scanner.Scan() {
		line := scanner.Text()
		for _, m := range metrics {
			m.MatchLine(line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	metfam, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		fmt.Println(err)
	}
	for _, met := range metfam {
		fmt.Println(met)
	}
}
