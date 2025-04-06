# Lexicon Beneficial Ownership Infrastructure

This repository contains the infrastructure configuration for the Lexicon Beneficial Ownership project. The project is designed to track and analyze beneficial ownership data through various microservices.

## Project Overview

The infrastructure is defined using Docker Compose and includes the following components:

- **Traefik**: Reverse proxy for routing requests
- **PostgreSQL**: Database for persistent storage
- **NATS**: Messaging system for service communication
- **Redis**: Caching layer
- **Multiple microservices**:
  - Beneficial Ownership API
  - Beneficial Ownership Frontend
  - Beneficial Ownership Dataminer
  - Beneficial Ownership Dashboard
  - Named Entity Recognition
  - Various court crawlers (Indonesia, Singapore)
  - HTTP Crawler Service

## Directory Structure

The project is organized as follows:

- `/v1`: Legacy version (contains docker-server)
- `/v2`: Current version with complete infrastructure
  - `docker-compose.yml`: Main service definitions
  - `.env`: Consolidated environment variables
  - `services-config.yaml`: Service configuration for environment management
  - Service-specific directories (each containing their own code)

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Make (optional, for using the Makefile shortcuts)

### Setup and Configuration

1. Clone this repository:
   ```
   git clone https://github.com/yourusername/lexicon-beneficial-ownership-infra.git
   cd lexicon-beneficial-ownership-infra/v2
   ```

2. Generate the consolidated environment file:
   ```
   ./deployment env
   ```

   Or using Make:
   ```
   make env
   ```

3. Update the Docker Compose configuration:
   ```
   ./deployment update
   ```

   Or using Make:
   ```
   make update-compose
   ```

4. Start the services:
   ```
   docker compose up -d
   ```

   Or using Make:
   ```
   make up
   ```

### Available Make Commands

- `make env`: Generate consolidated environment file
- `make update-compose`: Update docker-compose.yml with environment variables
- `make up`: Start all services
- `make down`: Stop all services

## Environment Configuration Management

The project uses a centralized environment variable management approach where all variables are defined in a single consolidated `.env` file. This makes it easier to maintain consistent configuration across all services.

### Service Configuration

Services are defined in the `services-config.yaml` file with the following structure:

```yaml
# Common infrastructure services
common_services:
  - name: postgres
    env_file: postgres/.env
    prefix: "POSTGRES_"
  # Other common services...

# Application-specific services
services:
  - name: lexicon-beneficial-ownership-api
    env_file: lexicon-beneficial-ownership-api/.env
    prefix: "BO_API_"
    domain: "beneficial-ownership.lexicon.id/api"
  # Other application services...
```

### Adding a New Service

To add a new service:

1. Create a directory for your service
2. Add service-specific `.env` file
3. Update the `services-config.yaml` with your service details
4. Run `./deployment env` to update the consolidated environment file
5. Run `./deployment update` to update the Docker Compose configuration

## Deployment Tools

The `deployment` tool provides several options:

### Environment Consolidation

```
./deployment env [options]
```

Options:
- `-o string`: Output file path (default: `.env`)
- `-f`: Force overwrite output file if it exists
- `-c string`: Path to services configuration file (default: `services-config.yaml`)
- `-d`: Auto-discover services in project directory
- `-dir string`: Directory to discover services (default: current directory)

### Docker Compose Update

```
./deployment update [options]
```

Options:
- `-t string`: Path to template file (default: `docker-compose.template.yml`)
- `-env string`: Path to consolidated env file (default: `.env`)
- `-o string`: Output file path (default: `docker-compose.yml`)
- `-f`: Force overwrite output file if it exists
- `-dir string`: Directory to discover services
- `-c string`: Path to services configuration file (default: `services-config.yaml`)

## License

MIT License

Copyright (c) 2023 Lexicon Beneficial Ownership

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## Contributors

- Adryan Eka Vandra ([@adryanev](https://github.com/adryanev))
- Muhammad Hanif Ramadhan ([@mhanifrmd](https://github.com/mhanifrmd))
- Naufal Fawaz Andriawan ([@andriawan24](https://github.com/andriawan24))