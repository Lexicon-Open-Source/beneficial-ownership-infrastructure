# Lexicon Beneficial Ownership Infrastructure

This repository contains the infrastructure configuration for the Lexicon Beneficial Ownership project.

## Overview

The infrastructure is defined using Docker Compose and includes the following services:
- Traefik (reverse proxy)
- PostgreSQL (database)
- NATS (messaging)
- Redis (caching)
- Multiple microservices (crawlers, APIs, frontend, dashboard)

## Environment Configuration Management

The project uses a centralized environment variable management approach where all variables are defined in a single consolidated `.env` file in this directory. This makes it easier to maintain consistent configuration across all services.

### Environment Management Scripts

This project includes several scripts to help manage environment variables:

1. **consolidate-env-files.sh**: Consolidates individual service .env files into a single main .env file
2. **update-docker-compose.sh**: Updates docker-compose.yml to use variables from the consolidated .env file
3. **run-env-workflow.sh**: Runs the full workflow (consolidate + update + optional start)

#### Consolidating Service .env Files

```bash
# Basic usage (uses defaults)
./consolidate-env-files.sh

# Auto-discover services with .env files
./consolidate-env-files.sh -d

# Use a custom configuration
./consolidate-env-files.sh -c services-config.yaml

# Specify output file
./consolidate-env-files.sh -o ./custom-env-path
```

#### YAML Configuration

The script uses a YAML configuration file (`services-config.yaml`) to define services and their prefixes:

```yaml
# Service prefixes
prefixes:
  postgres: "POSTGRES_"
  lexicon-beneficial-ownership-api: "BO_API_"
  # ... other service prefixes

# Services definitions
services:
  - name: postgres
    env_file: ../postgres/.env

  - name: lexicon-beneficial-ownership-api
    env_file: ../lexicon-beneficial-ownership-api/.env

  # ... other services
```

This YAML format makes it easy to add, modify, or remove services from the configuration.

#### Full Workflow

The `run-env-workflow.sh` script automates the entire environment management process:

```bash
# Basic usage
./run-env-workflow.sh

# Auto-discover services and start containers
./run-env-workflow.sh -d -s

# Create backups of files before modifying
./run-env-workflow.sh -b
```

This workflow script works on both macOS and Linux platforms.

### Consolidated Docker Compose Configuration

Docker Compose has been configured to use only the consolidated `.env` file rather than loading individual service-specific `.env` files. This approach:

1. Reduces complexity and duplication
2. Makes it easier to see all environment variables in one place
3. Ensures consistent variable values across services
4. Simplifies deployment and configuration

Each service in the docker-compose.yml has an `environment` section that maps variables from the consolidated `.env` file to the appropriate service-specific variables using Docker Compose variable substitution. For example:

```yaml
services:
  singapore-supreme-court-crawler:
    # ...
    environment:
      - APP_URL=${SINGAPORE_CRAWLER_APP_URL}
      - LISTEN_HOST=${SINGAPORE_CRAWLER_LISTEN_HOST}
      - POSTGRES_DB_NAME=${SINGAPORE_CRAWLER_POSTGRES_DB_NAME}
      # ...
    env_file:
      - .env
```

This approach maintains compatibility with local development while simplifying deployment.

### Dynamic Service Configuration

Services are defined in a YAML configuration file (`services-config.yaml`):

1. The `prefixes` section maps service names to their environment variable prefixes
2. The `services` section defines each service with:
   - Service name
   - Path to the service-specific .env file

This makes it easy to add or modify services without changing the core logic of the scripts.

### Adding a New Service

To add a new service to the environment variable system:

1. Create a service-specific .env file with the necessary variables
2. Add the service to the `services-config.yaml` file (if not using auto-discovery):
   ```yaml
   # Add prefix to the prefixes section
   prefixes:
     # ... existing prefixes
     new-service: "NEW_SERVICE_"

   # Add service to the services section
   services:
     # ... existing services
     - name: new-service
       env_file: ../new-service/.env
   ```
3. Run the consolidation workflow to update the main .env file
4. Add the service to `docker-compose.yml` with appropriate environment mappings:
   ```yaml
   service-name:
     # ...
     environment:
       - TARGET_VAR1=${NEW_SERVICE_VAR1}
       - TARGET_VAR2=${NEW_SERVICE_VAR2}
     env_file:
       - .env
   ```

### Adding or Modifying Environment Variables

When you need to add or modify environment variables:

1. Edit the service-specific .env file
2. Run the consolidation workflow to update the main .env file:
   ```bash
   ./run-env-workflow.sh -d
   ```

3. Or simply restart the entire stack with:
   ```bash
   ./run-env-workflow.sh -d -s
   ```

## Environment Configuration

This project uses a YAML configuration file for managing environment variables across multiple services.

### Configuration File

The environment configuration is managed in `services-config.yaml`, which contains three main sections:

1. `common_services`: Infrastructure services with shared environment variables (processed first)
2. `prefixes`: Maps service names to their environment variable prefixes
3. `services`: Defines application-specific services with their name and path to the `.env` file

Example configuration:

```yaml
# Common infrastructure services (processed first to avoid duplication)
common_services:
  - name: postgres
    env_file: ./postgres/.env
    prefix: "POSTGRES_"
  - name: nats
    env_file: ./nats/.env
    prefix: "NATS_"

# Service prefixes
prefixes:
  postgres: "POSTGRES_"
  nats: "NATS_"
  lexicon-beneficial-ownership-api: "BO_API_"

# Application services
services:
  - name: lexicon-beneficial-ownership-api
    env_file: ../lexicon-beneficial-ownership-api/.env
```

The `common_services` section is processed first, ensuring that shared infrastructure variables (like database credentials) aren't duplicated when multiple services reference them.

### Running the Environment Consolidation Workflow

The environment consolidation workflow combines individual service `.env` files into a single `.env` file with appropriate prefixes.

To run the workflow:

```bash
cd lexicon-beneficial-ownership-infra
./run-env-workflow.sh
```

The `consolidate-env-files.sh` script provides several options:

```bash
./consolidate-env-files.sh [OPTIONS]

Options:
  -h, --help              Display help message
  -o, --output <path>     Specify output file path (default: ./.env)
  -f, --force             Overwrite existing output file without prompting
  -c, --config <path>     Specify services configuration file path (YAML)
  -d, --discover          Auto-discover services by scanning directories
```

Examples:

```bash
# Use default settings (recommended)
./consolidate-env-files.sh

# Force overwrite existing .env file
./consolidate-env-files.sh --force

# Specify a different output location
./consolidate-env-files.sh --output /path/to/output/.env

# Use a different configuration file
./consolidate-env-files.sh --config ./my-custom-config.yaml

# Auto-discover services instead of using the config file
./consolidate-env-files.sh --discover
```

The script will combine all service `.env` files, prefixing variables according to the configuration.

### Adding a New Service

To add a new service:

1. Add the service prefix to the `prefixes` section of `services-config.yaml`:
   ```yaml
   prefixes:
     existing-service: "EXISTING_"
     new-service: "NEW_SERVICE_"
   ```

2. Add the service definition to the `services` section of `services-config.yaml`:
   ```yaml
   services:
     - name: existing-service
       env_file: ../existing-service/.env
     - name: new-service
       env_file: ../new-service/.env
   ```

3. Run the consolidation workflow to update the main `.env` file:
   ```bash
   ./run-env-workflow.sh
   ```

4. Add the service to `docker-compose.yml` with appropriate environment mappings:
   ```yaml
   services:
     new-service:
       image: new-service-image
       environment:
         - NEW_SERVICE_API_KEY=${NEW_SERVICE_API_KEY}
         - NEW_SERVICE_DB_URL=${NEW_SERVICE_DB_URL}
   ```

## Getting Started

1. Clone this repository
2. Create .env files for your services or use auto-discovery
3. Start the environment workflow:
   ```bash
   # Auto-discover services, create backups, and start containers
   ./run-env-workflow.sh -d -b -s
   ```

## Environment Management Approach

This project uses a bottom-up approach to environment management:

### Bottom-up Approach (Service-specific → Consolidated)

This approach is ideal when:
- Individual services are developed by different teams
- Each service has its own .env file with service-specific variables
- You need to create a consolidated environment for deployment

Steps:
1. Develop services with their own .env files
2. Run `consolidate-env-files.sh` to create a consolidated .env file
3. Run `update-docker-compose.sh` to update docker-compose.yml
4. Deploy using the consolidated environment

## Service Dependencies

The infrastructure services have the following dependencies:
- All services depend on the `setup` service that runs the environment workflow
- API and frontend services depend on PostgreSQL, NATS, and Redis
- Crawler services depend on PostgreSQL and NATS
