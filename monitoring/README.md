# Monitoring Stack (LGTM)

This directory contains the configuration for the LGTM (Loki, Grafana, Tempo, Mimir) monitoring stack.

## Components

- **Loki**: Log aggregation system
- **Grafana**: Visualization and dashboards
- **Tempo**: Distributed tracing backend
- **Mimir**: Metrics platform for long-term storage
- **Prometheus**: Metrics collection and alerting

## Configuration

All configuration is managed through:

1. YAML configuration files for each service
2. Environment variables in the `.env` file
3. Docker Swarm stack definition in `monitoring-stack.yml`

## Environment Variables

The `.env` file contains all configurable parameters:

- Grafana admin credentials
- Resource limits for each service
- Health check parameters
- Service ports
- Domain configuration

To set up your environment:

```bash
# Copy the example environment file
cp .env.example .env

# Edit the .env file with your specific configuration
nano .env
```

## Deployment

To deploy the monitoring stack:

```bash
# Make sure you're in the monitoring directory
cd monitoring

# Deploy the stack
docker stack deploy -c monitoring-stack.yml monitoring
```

## Access

Grafana dashboard is accessible at:
- http://monitoring.beneficialowner.lexicon.id (via Traefik)
- http://localhost:3000 (direct access)

Default credentials:
- Username: admin
- Password: admin (change this in production)

## Data Sources

The following data sources are automatically configured:

- Prometheus: Metrics from the application
- Loki: Logs from all services
- Tempo: Distributed traces
- Mimir: Long-term metrics storage