package otelx

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.38.0"
	"google.golang.org/grpc"
)

type ShutdownFn func(context.Context) error

type Monitor struct {
	serviceName string
	grpcConn    *grpc.ClientConn
}

func NewMonitor(serviceName string, conn *grpc.ClientConn) *Monitor {
	return &Monitor{
		serviceName: serviceName,
		grpcConn:    conn,
	}
}

func (m *Monitor) Close() error {
	return m.grpcConn.Close()
}

func (m *Monitor) InitMeterProvider(ctx context.Context) (ShutdownFn, error) {
	res, err := resource.New(ctx,
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(m.serviceName),
		))
	if err != nil {
		return nil, err
	}

	otlpExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithGRPCConn(m.grpcConn),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	options := []sdkmetric.Option{
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(otlpExporter, sdkmetric.WithInterval(5*time.Second))),
	}

	promExporter, err := newPrometheusExporter()
	if err != nil {
		return nil, err
	}

	options = append(options, sdkmetric.WithReader(promExporter))
	meterProvider := sdkmetric.NewMeterProvider(options...)

	otel.SetMeterProvider(meterProvider)
	return meterProvider.Shutdown, nil
}

func (m *Monitor) InitTracerProvider(ctx context.Context) (ShutdownFn, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(m.serviceName),
		))
	if err != nil {
		return nil, err
	}

	otlpExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithGRPCConn(m.grpcConn),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(otlpExporter, sdktrace.WithBatchTimeout(5*time.Second)),
	)

	otel.SetTracerProvider(tracerProvider)
	return tracerProvider.Shutdown, nil
}
