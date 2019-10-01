package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/resty.v1"

	log "github.com/sirupsen/logrus"
)

type tokenQuery struct {
	ApiToken string `json:"apiToken"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

type mpulseHistogram struct {
	ChartTitle       string                `json:"chartTitle"`
	ChartTitleSuffix string                `json:"chartTitleSuffix"`
	DatasetName      string                `json:"datasetName"`
	ReportType       string                `json:"reportType"`
	ResultName       string                `json:"resultName"`
	Series           mpulseHistogramSeries `json:"series"`
}

type mpulseHistogramSeries struct {
	Series []mpulseHistogramSeriesSeries `json:"series"`
}

type mpulseHistogramSeriesSeries struct {
	Name           string `json:"name"`
	KValue         int32  `json:"kValue"`
	Median         int32  `json:"median"`
	PercentileName string `json:"percentile_name"`
	P95            int32  `json:"p95"`
	P98            int32  `json:"p98"`
	Buckets        int32  `json:"buckets"`
}

type mpulseCollector struct {
	timerMetric *prometheus.Desc
	host        string
	key         string
	timers      []string
	token       string
}

func newMpulseCollector(host string, key string, token string, timers []string) *mpulseCollector {
	log.Info("Obtaining new security token from mPulse...")

	client := resty.New()
	path := "services/rest/RepositoryService/v1/Tokens"
	resp, err := client.R().
		SetBody(tokenQuery{ApiToken: token}).
		SetResult(&tokenResponse{}).
		Put(fmt.Sprintf("%s/%s", host, path))

	if err != nil {
		log.Fatal(err)
	}

	result := resp.Result().(*tokenResponse)
	securityToken := result.Token

	if securityToken != "" {
		log.Info("Successfully obtained security token")
	} else {
		log.Fatal("Error obtaining security token!")
	}

	return &mpulseCollector{
		timerMetric: prometheus.NewDesc("mpulse_timer_metric",
			"Information about mPulse timer percentiles",
			[]string{"timer", "quantile"}, nil),
		host:   host,
		key:    key,
		timers: timers,
		token:  securityToken,
	}
}

func (collector *mpulseCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.timerMetric
}

func (collector *mpulseCollector) Collect(ch chan<- prometheus.Metric) {
	url := fmt.Sprintf(
		"%s/mpulse/api/v2/%s/histogram",
		collector.host,
		collector.key,
	)

	client := resty.New()
	client.SetHeader("Authentication", collector.token)

	for _, timer := range collector.timers {
		resp, err := client.R().
			SetQueryParams(map[string]string{
				"timer":            timer,
				"date-comparator":  "Last",
				"trailing-seconds": "20",
				"series-format":    "json",
			}).
			SetResult(&mpulseHistogram{}).
			Get(url)

		if err != nil {
			log.Fatal(err)
		}

		histogram := resp.Result().(*mpulseHistogram)
		series := histogram.Series.Series[0]

		ch <- prometheus.MustNewConstMetric(
			collector.timerMetric, prometheus.GaugeValue, (float64(series.Median) / 1000), timer, "0.5")
		ch <- prometheus.MustNewConstMetric(
			collector.timerMetric, prometheus.GaugeValue, (float64(series.P95) / 1000), timer, "0.95")
		ch <- prometheus.MustNewConstMetric(
			collector.timerMetric, prometheus.GaugeValue, (float64(series.P98) / 1000), timer, "0.98")
	}
}
