package main

import (
	"flag"
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	log "github.com/sirupsen/logrus"
)

func getEnv(key string, defaultVal string) string {
	if envVal, ok := os.LookupEnv(key); ok {
		return envVal
	}
	return defaultVal
}

func main() {
	var (
		apiKey = flag.String("api.key", getEnv("MPULSE_API_KEY", ""),
			"API Key of mPulse application")
		apiToken = flag.String("api.token", getEnv("MPULSE_API_TOKEN", ""),
			"User's API token")
		mpulseTimers = flag.String("histogram.timers", getEnv("MPULSE_HISTOGRAM_TIMERS", ""),
			"Histogram timers to scrape")
		logFormat = flag.String("log-format", getEnv("MPULSE_EXPORTER_LOG_FORMAT", "txt"),
			"Log format, valid options are txt and json")
	)
	flag.Parse()

	switch *logFormat {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.SetFormatter(&log.TextFormatter{})
	}

	log.Info("mPulse Metrics Exporter")

	mpulse := newMpulseCollector("https://mpulse.soasta.com/concerto",
		*apiKey, *apiToken, strings.Split(*mpulseTimers, ","))
	prometheus.MustRegister(mpulse)

	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve metrics on :8080/metrics")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
