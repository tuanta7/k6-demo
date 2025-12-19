# SigNoz

## DurationNano

Every span in OTLP has built-in timestamps `start_time_unix_nano` and `end_time_unix_nano`

- SigNoz stores spans in ClickHouse, and for convenience it materializes that duration into a field often called `durationNano`
- In the UI it may display the value in milliseconds for readability.