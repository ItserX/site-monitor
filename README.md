# Website Monitoring System

## Project Overview

A distributed microservices-based system for monitoring website availability, collecting performance metrics, and sending real-time alerts. Built with modern cloud-native technologies and designed for scalability and reliability.

## System Architecture

### Core Services

| Service | Description | Port | Database |
|---------|-------------|------|----------|
| **CRUD Service** | Manages website entities in database | 8080 | PostgreSQL |
| **Checker Service** | Performs website availability checks | - | - |
| **Alert Service** | Alerts delivery | - | - |

### Monitoring Stack

| Component | Purpose | Port |
|-----------|---------|------|
| **Prometheus** | Metrics collection and storage | 9090 |
| **Loki** | Log aggregation and storage | 3100 |
| **Grafana** | Visualization and dashboards | 3000 |
| **Promtail** | Log collection agent | - |
| **Pushgateway** | Metrics gateway | 9091 |

### Data Storage

| Database | Purpose |
|----------|---------|
| **PostgreSQL** | Persistent storage of website configurations |
| **Redis** | Temporary storage of check statuses and alerts |

## Project Structure

```bash
.
├── cmd/                    # Application entry points
│   ├── alert/             # Alert service main
│   ├── checker/           # Checker service main
│   └── crud/              # CRUD service main
├── configs/               # Service configuration files
│   ├── alert.yaml
│   ├── checker.yaml
│   └── crud.yaml
├── internal/              # Internal application code
│   ├── alert/             # Alert business logic
│   ├── checker/           # Checker business logic
│   ├── config/            # Configuration management
│   ├── crud/              # HTTP handlers
│   ├── storage/           # Database abstractions
│   └── telegram/          # Telegram integration
├── pkg/                   # Shared utilities
│   ├── logger/            # Structured logging
│   ├── metrics/           # Prometheus metrics
│   └── utils/             # Common utilities
├── migrations/            # Database schema migrations
├── grafana/               # Grafana provisioning
│   └── provisioning/
│       ├── dashboards/    # Dashboard definitions
│       └── datasources/   # Data source configurations
├── docker-compose.yaml    # Full environment setup
├── Dockerfile.*           # Service-specific Dockerfiles
└── *.yml                  # Monitoring configuration files
```

## API
### CRUD Service (8080)
```bash 
GET    /sites          # List all websites
GET    /sites/{id}     # Get specific website
POST   /sites          # Add new website
PUT    /sites/{id}     # Update website
DELETE /sites/{id}     # Remove website
```

## Deployment
```bash 
cd site-monitor
docker-compose up
```
