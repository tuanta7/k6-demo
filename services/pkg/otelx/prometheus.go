package otelx

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
)

func newPrometheusExporter() (*otelprom.Exporter, error) {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewBuildInfoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	promExporter, err := otelprom.New(
		otelprom.WithRegisterer(registry),
	)
	if err != nil {
		return nil, err
	}

	return promExporter, nil
}
