# Asynchronous Event & Notification Service
Go microservice using Redis BRPOP queues, goroutines, channels, a bounded worker pool, retries/dead-letter handling, graceful shutdown, Docker, and CI.

Run: `go test -race ./...` and `docker compose up --build`.

Synthetic concurrency tests are included; no specific production throughput is claimed.
