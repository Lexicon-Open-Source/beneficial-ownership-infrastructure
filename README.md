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

### Environment Management Tools

This project includes deployment tools to help manage environment variables:

1. **`deployment env`**: Consolidates individual service .env files into a single main .env file
2. **`deployment update`**: Updates docker-compose.yml to use variables from the consolidated .env file

#### Consolidating Service .env Files

```bash
# Basic usage (uses defaults)
deployment env

# Auto-discover services with .env files
deployment env -d

# Use a custom configuration
deployment env -c services-config.yaml

# Specify output file
deployment env -o ./custom-env-path

# Specify a custom directory to discover services
deployment env -d -dir ./services
```

#### YAML Configuration

The tool uses a YAML configuration file (`services-config.yaml`) to define services and their prefixes:

```yaml
# Common infrastructure services (processed first to avoid duplication)
common_services:
  - name: postgres
    env_file: postgres/.env
    prefix: "POSTGRES_"
  - name: nats
    env_file: nats/.env
    prefix: "NATS_"

# Services definitions
services:
  - name: lexicon-beneficial-ownership-api
    env_file: lexicon-beneficial-ownership-api/.env
    prefix: "BO_API_"

  # ... other services
```

This YAML format makes it easy to add, modify, or remove services from the configuration.

#### Updating Docker Compose

The `deployment update` command updates the docker-compose.yml file to use the consolidated environment variables:

```bash
# Basic usage
deployment update

# Specify custom paths
deployment update -env .env -o docker-compose.yml

# Use a custom template
deployment update -t docker-compose.template.yml

# Specify a custom directory to discover services
deployment update -dir ./services
```

### Dynamic Service Configuration

Services are defined in a YAML configuration file (`services-config.yaml`):

1. The `common_services` section defines shared infrastructure services
2. The `services` section defines application-specific services with:
   - Service name
   - Path to the service-specific .env file
   - Environment variable prefix

Each service in the docker-compose.yml has its environment variables configured to reference the consolidated .env file.

### Adding a New Service

To add a new service to the environment variable system:

1. Create a service-specific .env file with the necessary variables
2. Add the service to the `services-config.yaml` file (if not using auto-discovery):
   ```yaml
   # Add service to the services section
   services:
     # ... existing services
     - name: new-service
       env_file: new-service/.env
       prefix: "NEW_SERVICE_"
   ```
3. Run the consolidation command to update the main .env file:
   ```bash
   deployment env
   ```
4. Update docker-compose.yml to include the new service:
   ```bash
   deployment update
   ```

### Customizing Service Discovery

The deployment tools support customizing where to look for services:

1. Using the `-dir` parameter to specify a custom directory:
   ```bash
   # Look for services in the ./services directory
   deployment env -d -dir ./services

   # Update docker-compose using services in the ./services directory
   deployment update -dir ./services
   ```

2. This is useful when:
   - Your services are organized in a different directory structure
   - You have multiple environments with services in different locations
   - You want to generate configurations for a subset of services

### Environment Configuration

This project uses a YAML configuration file for managing environment variables across multiple services.

### Configuration File

The environment configuration is managed in `services-config.yaml`, which contains two main sections:

1. `common_services`: Infrastructure services with shared environment variables (processed first)
2. `services`: Defines application-specific services with their name, env_file path, and prefix

Example configuration:

```yaml
# Common infrastructure services (processed first to avoid duplication)
common_services:
  - name: postgres
    env_file: postgres/.env
    prefix: "POSTGRES_"
  - name: nats
    env_file: nats/.env
    prefix: "NATS_"

# Application services
services:
  - name: lexicon-beneficial-ownership-api
    env_file: lexicon-beneficial-ownership-api/.env
    prefix: "BO_API_"
```

The `common_services` section is processed first, ensuring that shared infrastructure variables (like database credentials) aren't duplicated when multiple services reference them.

### Running the Environment Consolidation Command

The environment consolidation command combines individual service `.env` files into a single `.env` file with appropriate prefixes.

To run the command:

```bash
cd lexicon-beneficial-ownership-infra
deployment env
```

The `deployment env` command provides several options:

```
Options:
  -o string   Output file path for consolidated env file (default: .env)
  -f          Force overwrite output file if it exists
  -c string   Path to services configuration file (default: services-config.yaml)
  -d          Auto-discover services in project directory
  -dir string Directory to discover services (default: current directory)
                Use this to specify a different directory for service discovery
```

Examples:

```bash
# Use default settings
deployment env

# Force overwrite existing .env file
deployment env -f

# Specify a different output location
deployment env -o /path/to/output/.env

# Use a different configuration file
deployment env -c ./my-custom-config.yaml

# Auto-discover services in a specific directory
deployment env -d -dir ./services
```

### Updating Docker Compose Configuration

To update Docker Compose to use the consolidated environment variables:

```bash
deployment update
```

The `deployment update` command provides several options:

```
Options:
  -t string     Path to template file (default: docker-compose.template.yml in project root)
  -env string   Path to consolidated env file (default: .env)
  -o string     Output file path (default: docker-compose.yml in project root)
  -f            Force overwrite output file if it exists
  -dir string   Directory to discover services (default: current directory)
  -c string     Path to services configuration file (default: services-config.yaml)
```

Examples:

```bash
# Use default settings
deployment update

# Specify custom input and output files
deployment update -t my-template.yml -env ./my-env-file -o ./my-docker-compose.yml

# Auto-discover services in a specific directory
deployment update -dir ./services
```

### Adding a New Service

To add a new service to the environment variable system:

1. Create a service-specific .env file with the necessary variables
2. Add the service to the `services-config.yaml` file (if not using auto-discovery):
   ```yaml
   # Add service to the services section
   services:
     # ... existing services
     - name: new-service
       env_file: new-service/.env
       prefix: "NEW_SERVICE_"
   ```
3. Run the consolidation command to update the main .env file:
   ```bash
   deployment env
   ```
4. Update docker-compose.yml to include the new service:
   ```bash
   deployment update
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
