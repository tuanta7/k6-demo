# Monitoring

## 1. Fundamentals

The three pillars of observability are metrics, logs, and traces.

### Minimal Setup

| Component               | Role                                                                                                  |
| ----------------------- | ----------------------------------------------------------------------------------------------------- |
| OpenTelemetry SDKs      | Instrument services with automatic or manual traces, logs and metrics.                                |
| OpenTelemetry Collector | Collects metrics, traces, logs; forwards to SigNoz Collector                                          |
| SigNoz (All-in-one)     | Observability platform native to OpenTelemetry with logs, traces and metrics in a single application. |
| ClickHouse              | A lightning-fast, column-oriented database engineered specifically for analytical workloads           |

## 2. OpenTelemetry

A major goal of OpenTelemetry is to enable easy instrumentation of any applications and systems, regardless of the programming language, infrastructure, and runtime environments used. It can be used to instrument, generate, collect, and export all 3 types of telemetry data.

## 3. Profiling

Profiling is a technique used in software development to measure and analyze the runtime behavior of a program.

Profilers, particularly in languages like Go, observe compiled and running code at the binary or runtime level. They determine which function based on symbols (function names and locations) and stack traces, which are available at runtime or via debug information.

> [!NOTE]
> In interpreted or runtime-managed languages like JavaScript, Python, and Java, profilers operate quite differently from compiled languages like Go or C++. These languages provide rich runtime introspection features, which allow profilers to obtain function names, stack traces, and timing information without relying on binary symbol tables.

### Profiling vs Tracing

- Tracing shows where a request started, which rservices/functions it took, and how long each segment took. Each trace is tied to a specific request.
- Profiling is not about tracking requests at all; it's about understanding where the application spends compute or memory resources, often regardless of which request triggered them.
