# BookCommunity

A high-performance book community backend built with Go.

## Tech Stack

- **Backend**: Go 1.20 + Gin + GORM
- **Database**: PostgreSQL 15
- **Cache**: Redis 7.0 + In-Memory LRU
- **Message Queue**: RabbitMQ 3.12
- **Monitoring**: Prometheus + Grafana
- **Deployment**: Docker + Kubernetes

## Quick Start

```bash
# Clone the repository
git clone https://github.com/sylvia-ymlin/Coconut-book-community.git
cd Coconut-book-community

# Start services with Docker Compose
docker-compose up -d

# Copy configuration
cp config/conf/example.yaml config/conf/config.yaml

# Run the application
go run main.go
```

## Architecture

```
┌─────────────────────────────────┐
│       BookCommunity API         │
│    Go + Gin + GORM              │
└──────────┬──────────────────────┘
           │
    ┌──────┼──────┐
    ↓      ↓      ↓
┌────────┐ │  ┌────────┐
│Postgres│ │  │RabbitMQ│
└────────┘ │  └────────┘
           ↓
      ┌────────┐
      │ Redis  │
      └────────┘
```

## Performance

- **QPS**: 5000+ requests/sec
- **Latency**: P99 < 50ms
- **Cache Hit Rate**: 95%+

## License

MIT
