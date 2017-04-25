package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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
	totalLines = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "stream_exporter",
			Subsystem: "line_processing",
			Name:      "lines_total",
			Help:      "Number of lines processed",
		},
	)
)

var (
	configFilePath = flag.String("config", "stream_exporter.yaml", "Path to config file")

	inputType      = flag.String("input.type", "", "What input module to use")
	listInputTypes = flag.Bool("input.print", false, "Print available input modules and exit")

	metricsListenAddr = flag.String("web.listen-address", ":9177", "Address on which to expose metrics")
	metricsPath       = flag.String("web.metrics-path", "/metrics", "Path under which the metrics are available")
)

func main() {
	flag.Parse()

	if *listInputTypes {
		fmt.Println(input.GetAvailableInputs())
		os.Exit(0)
	}
	if *inputType == "" {
		fmt.Printf("-input.type is required. The following input types are available:\n%v", input.GetAvailableInputs())
		os.Exit(1)
	}

	metricsConfig, err := linemetrics.ReadPatternConfig(*configFilePath)
	if err != nil {
		log.Fatalf("Could not read pattern config: %v", err)
	}

	// Define metrics
	metrics := make([]linemetrics.LineMetric, 0, len(metricsConfig))
	for _, definition := range metricsConfig {
		lineMetric, collector := linemetrics.NewLineMetric(definition)
		metrics = append(metrics, lineMetric)
		prometheus.MustRegister(collector)
	}

	prometheus.MustRegister(lineProcessingTime)
	prometheus.MustRegister(totalLines)

	// Setup signal handling
	quitSig := make(chan os.Signal, 1)
	signal.Notify(quitSig, os.Interrupt)

	// Configure input
	inputReader, err := input.NewInput(*inputType)
	if err != nil {
		log.Fatalf("Could not initialize input: %v", err)
	}

	inputChannel := make(chan string)
	go inputReader.StartStream(inputChannel)

	// Setup http server
	http.Handle(*metricsPath, promhttp.Handler())
	go http.ListenAndServe(*metricsListenAddr, nil)

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
			totalLines.Inc()
		case <-quitSig:
			log.Info("Received quit signal, shutting down...")
			done = true
			break
		}
	}
}
