# Akamai mPulse Prometheus exporter

Prometheus exporter for Akamai mPulse metrics.

Support histogram metrics of timers

## Running
### Flags

Name                   | Environment Variable Name            | Description
-----------------------|--------------------------------------|-----------------
api.key                | MPULSE_API_KEY                       | API Key of mPulse application
api.token              | MPULSE_API_TOKEN                     | User's API token
histogram.timers       | MPULSE_HISTOGRAM_TIMERS              | Histogram timers to scrape
log-format             | MPULSE_EXPORTER_LOG_FORMAT           | Log format, valid options are txt and json
