# Monorepo for Location Services

This repository will grow gradually to keep the codebase simple; early package splitting will be avoided whenever
possible.

## Lesson Learned

- Accept interfaces, return structs.

### Adapters

- Database and queue adapters (wrappers/clients) should be implemented in the `pkg` directory.
- Define the adapter interfaces where they are consumed. These adapter interfaces should abstract driver or library
  differences for the same backend (for example, switching from `lib/pq` to `pgx`, or `go-redis` to `redigo`), 
  not high‑level changes such as PostgreSQL to MongoDB or RabbitMQ to Kafka.
- Sometimes a library can be trusted to be never deprecated, so it can be used directly without mapping all of its
  types to a local one.
- The interfaces that abstract the backend switching (such as PostgreSQL to MongoDB) should be defined in the higher
  level, like the `repository` layer.