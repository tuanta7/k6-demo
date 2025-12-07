# Metrics in Go

Metrics are numerical time-series data collection that track the health and performance of a system over time. They help identify trends, anomalies, and overall system health. They can represent:

- **Infrastructure-level metrics** (often collected by agents): CPU, Memory, Disk I/O, Network throughput, etc
- **Application-level metrics**: Request counts per endpoint, Latency histograms, Cache hit/miss ratio, Queue length, etc
- **Custom business metrics**: Number of orders placed, Items in shopping cart, Messages processed per second, etc.

## 1. Prometheus

Prometheus supports instrumentation libraries for multiple programming languages, including [Go](https://github.com/prometheus/client_golang). The Go client library has two separate parts, one for instrumenting application code, and one for creating clients that talk to the Prometheus HTTP API.

Expose the default metrics via an HTTP endpoint. The endpoint must be named `/metrics`

```go
import (
    "net/http"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
    http.Handle("/metrics", promhttp.Handler())
    http.ListenAndServe(":2112", nil)
}
```

Configure Prometheus server to scrape the endpoint every 10s. Metric values are kept in memory inside the application until scraped.

- An endpoint that can be scraped is called an _INSTANCE_, usually corresponding to a single process.
- A collection of instances with the same purpose, a process replicated for scalability or reliability for example, is called a _JOB_.

```yml
scrape_configs:
  - job_name: myapp
    scrape_interval: 10s
    static_configs:
      - targets:
          - localhost:2112
```

After metrics are collected from instrumented applications, the Prometheus server functions as a time-series database which Grafana can query directly using Prometheus's PromQL to render dashboards, graphs, and alerts.

> [!NOTE]
> However, Prometheus is not a general-purpose database—it's optimized for time-series only and lacks traditional features like joins or full-text search (those are handled by log stores like Loki or Elasticsearch).

## 2. OpenTelemetry

OpenTelemetry can collect metrics using a consolidated API via push or pull, potentially transforms them, and sends them onward to other systems (often Prometheus or a Prometheus-compatible one like Grafana Mimir) for storage or query.

### Metrics API

The Metrics API consists of these main components:

- **MeterProvider**: the entry point of the API, holds meter configuration and provides access to Meters.
- **Meter**: responsible for creating Instruments.
- **Instrument**: responsible for reporting Measurements.
