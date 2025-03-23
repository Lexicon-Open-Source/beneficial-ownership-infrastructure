#!/bin/bash

# This script updates the environment variables in docker-compose.yml based on the consolidated .env file

# Define help function
show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo
    echo "Options:"
    echo "  -h, --help              Display this help message"
    echo "  -e, --env <path>        Specify the consolidated .env file path (default: ./.env)"
    echo "  -d, --docker-compose <path>  Specify the docker-compose.yml file path (default: ./docker-compose.yml)"
    echo "  -c, --config <path>     Specify services configuration file path (YAML)"
    echo "  -b, --backup            Create a backup of the docker-compose.yml file before updating"
    echo "  -f, --force             Force operations without prompting for confirmation"
    echo
    echo "This script updates docker-compose.yml environment variables based on a consolidated .env file."
}

# Check if yq is installed
check_dependencies() {
    if ! command -v yq &> /dev/null; then
        echo "Error: This script requires yq to parse YAML files."
        echo "Please install yq using one of the following methods:"
        echo "  - Homebrew (macOS): brew install yq"
        echo "  - Apt (Debian/Ubuntu): apt-get install yq"
        echo "  - Snap: snap install yq"
        echo "  - Go: go install github.com/mikefarah/yq/v4@latest"
        echo "  - Download binary from: https://github.com/mikefarah/yq/releases"
        exit 1
    fi
}

# Process command line arguments
ENV_FILE=""
DOCKER_COMPOSE_FILE=""
CONFIG_FILE=""
CREATE_BACKUP=false
FORCE_OPERATIONS=false

while [[ $# -gt 0 ]]; do
    case "$1" in
        -h|--help)
            show_help
            exit 0
            ;;
        -e|--env)
            if [[ -n "$2" && "$2" != -* ]]; then
                ENV_FILE="$2"
                shift 2
            else
                echo "Error: Argument for $1 is missing" >&2
                show_help
                exit 1
            fi
            ;;
        -d|--docker-compose)
            if [[ -n "$2" && "$2" != -* ]]; then
                DOCKER_COMPOSE_FILE="$2"
                shift 2
            else
                echo "Error: Argument for $1 is missing" >&2
                show_help
                exit 1
            fi
            ;;
        -c|--config)
            if [[ -n "$2" && "$2" != -* ]]; then
                CONFIG_FILE="$2"
                shift 2
            else
                echo "Error: Argument for $1 is missing" >&2
                show_help
                exit 1
            fi
            ;;
        -b|--backup)
            CREATE_BACKUP=true
            shift
            ;;
        -f|--force)
            FORCE_OPERATIONS=true
            shift
            ;;
        *)
            echo "Unknown option: $1" >&2
            show_help
            exit 1
            ;;
    esac
done

# Check dependencies
check_dependencies

INFRA_DIR=$(dirname "$0")

# Set default file paths if not specified
if [[ -z "$ENV_FILE" ]]; then
    ENV_FILE="$INFRA_DIR/.env"
fi

if [[ -z "$DOCKER_COMPOSE_FILE" ]]; then
    DOCKER_COMPOSE_FILE="$INFRA_DIR/docker-compose.yml"
fi

if [[ -z "$CONFIG_FILE" ]]; then
    CONFIG_FILE="$INFRA_DIR/services-config.yaml"
fi

# Check if files exist
if [ ! -f "$ENV_FILE" ]; then
    echo "Error: Consolidated .env file not found at $ENV_FILE"
    exit 1
fi

if [ ! -f "$DOCKER_COMPOSE_FILE" ]; then
    echo "Error: docker-compose.yml file not found at $DOCKER_COMPOSE_FILE"
    exit 1
fi

if [ ! -f "$CONFIG_FILE" ]; then
    echo "Warning: Services configuration file not found at $CONFIG_FILE"
    echo "The script will continue but may not correctly map service prefixes."
fi

# Prompt for confirmation if not forced
if [ "$FORCE_OPERATIONS" = false ]; then
    read -p "This will update the environment variables in $DOCKER_COMPOSE_FILE. Continue? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Operation cancelled."
        exit 0
    fi
fi

# Create backup if requested
if [ "$CREATE_BACKUP" = true ]; then
    BACKUP_FILE="${DOCKER_COMPOSE_FILE}.backup.$(date +%Y%m%d_%H%M%S)"
    echo "Creating backup of docker-compose.yml at $BACKUP_FILE"
    cp "$DOCKER_COMPOSE_FILE" "$BACKUP_FILE"
fi

echo "Updating docker-compose.yml with environment variables from $ENV_FILE"

# Load service prefixes from the configuration file
# Use regular variables instead of associative arrays for compatibility
declare -a PREFIX_KEYS
declare -a PREFIX_VALUES

if [ -f "$CONFIG_FILE" ]; then
    echo "Loading service prefixes from $CONFIG_FILE"

    # Load common service prefixes
    common_count=$(yq eval '.common_services | length' "$CONFIG_FILE")
    for ((i=0; i<$common_count; i++)); do
        service_name=$(yq eval ".common_services[$i].name" "$CONFIG_FILE")
        prefix=$(yq eval ".common_services[$i].prefix" "$CONFIG_FILE")
        PREFIX_KEYS+=("$service_name")
        PREFIX_VALUES+=("$prefix")
        echo "  - Common service: $service_name -> prefix: $prefix"
    done

    # Load app service prefixes
    service_count=$(yq eval '.services | length' "$CONFIG_FILE")
    for ((i=0; i<$service_count; i++)); do
        service_name=$(yq eval ".services[$i].name" "$CONFIG_FILE")
        prefix=$(yq eval ".services[$i].prefix" "$CONFIG_FILE")
        PREFIX_KEYS+=("$service_name")
        PREFIX_VALUES+=("$prefix")
        echo "  - App service: $service_name -> prefix: $prefix"
    done
fi

# Create a temporary file for processing
TEMP_FILE=$(mktemp)

# Function to get service name from indentation level
get_current_service() {
    local line=$1
    if [[ "$line" =~ ^[[:space:]]{2}([[:alnum:]-]+): ]]; then
        echo "${BASH_REMATCH[1]}"
    else
        echo "$current_service"
    fi
}

# Function to find the appropriate env var for a service variable
find_env_var() {
    local service=$1
    local var_name=$2

    # If this service has a prefix defined, try that first
    for ((i=0; i<${#PREFIX_KEYS[@]}; i++)); do
        if [[ "${PREFIX_KEYS[$i]}" == "$service" ]]; then
            echo "${PREFIX_VALUES[$i]}$var_name"
            return 0
        fi
    done

    # Try common services
    for common_service in postgres nats redis traefik; do
        for ((i=0; i<${#PREFIX_KEYS[@]}; i++)); do
            if [[ "${PREFIX_KEYS[$i]}" == "$common_service" ]]; then
                # Check if this variable exists in the .env file with the common service prefix
                if grep -q "^${PREFIX_VALUES[$i]}$var_name=" "$ENV_FILE"; then
                    echo "${PREFIX_VALUES[$i]}$var_name"
                    return 0
                fi
            fi
        done
    done

    # Fall back to the unprefixed variable name
    echo "$var_name"
    return 0
}

# Process the docker-compose.yml file
current_service=""
in_environment_section=false

# Read docker-compose.yml line by line
while IFS= read -r line || [[ -n "$line" ]]; do
    # Check if we're in a service definition
    service_check=$(get_current_service "$line")
    if [[ "$service_check" != "$current_service" ]]; then
        current_service="$service_check"
        in_environment_section=false
        echo "Processing service: $current_service"
    fi

    # Check if we're entering the environment section of a service
    if [[ "$line" =~ ^[[:space:]]{4}environment: ]]; then
        in_environment_section=true
        echo "$line" >> "$TEMP_FILE"
        continue
    fi

    # If we're in the environment section, process environment variables
    if [ "$in_environment_section" = true ]; then
        # Check if we're exiting the environment section
        if [[ ! "$line" =~ ^[[:space:]]{6}- ]]; then
            in_environment_section=false
            echo "$line" >> "$TEMP_FILE"
        else
            # This is an environment variable line
            if [[ "$line" =~ ^([[:space:]]{6}-[[:space:]])([A-Za-z0-9_]+)(\$\{[A-Za-z0-9_]+\}) ]]; then
                # This is a variable with a value using ${VAR} syntax
                spaces="${BASH_REMATCH[1]}"
                var_name="${BASH_REMATCH[2]}"

                # Find the appropriate env var reference
                env_var=$(find_env_var "$current_service" "$var_name")

                # Replace the line with the updated variable reference
                echo "$spaces$var_name=\${$env_var}" >> "$TEMP_FILE"
            elif [[ "$line" =~ ^([[:space:]]{6}-[[:space:]])([A-Za-z0-9_]+)=(.*) ]]; then
                # This is a variable with explicit value (VAR=value)
                spaces="${BASH_REMATCH[1]}"
                var_name="${BASH_REMATCH[2]}"
                value="${BASH_REMATCH[3]}"

                # If the value is a variable reference ${...}, update it
                if [[ "$value" =~ ^\$\{([A-Za-z0-9_]+)\} ]]; then
                    old_var="${BASH_REMATCH[1]}"

                    # Find the appropriate env var reference
                    env_var=$(find_env_var "$current_service" "$var_name")

                    # Replace the line with the updated variable reference
                    echo "$spaces$var_name=\${$env_var}" >> "$TEMP_FILE"
                else
                    # Keep the original line for non-variable values
                    echo "$line" >> "$TEMP_FILE"
                fi
            else
                # Keep other environment lines as is
                echo "$line" >> "$TEMP_FILE"
            fi
        fi
    else
        # Write all other lines to the temp file
        echo "$line" >> "$TEMP_FILE"
    fi
done < "$DOCKER_COMPOSE_FILE"

# Replace the original file with the temp file
mv "$TEMP_FILE" "$DOCKER_COMPOSE_FILE"

echo "docker-compose.yml has been updated successfully."