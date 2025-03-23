#!/bin/bash

# This script runs the full environment workflow:
# 1. Consolidate service-specific .env files into a single main .env file
# 2. Update docker-compose.yml to use the consolidated environment variables
# 3. (Optionally) start the services with docker-compose

# Define help function
show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo
    echo "Options:"
    echo "  -h, --help              Display this help message"
    echo "  -c, --config <path>     Specify services configuration file path (YAML)"
    echo "  -d, --discover          Auto-discover services by scanning directories"
    echo "  -s, --start             Start services after updating environment"
    echo "  -f, --force             Force overwrite of existing files without prompting"
    echo "  -b, --backup            Create backups of important files before modifying"
    echo
    echo "This script runs the full environment workflow for your microservices project."
}

# Process command line arguments
CONFIG_FILE=""
AUTO_DISCOVER=false
START_SERVICES=false
FORCE_OVERWRITE=false
CREATE_BACKUP=false

while [[ $# -gt 0 ]]; do
    case "$1" in
        -h|--help)
            show_help
            exit 0
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
        -d|--discover)
            AUTO_DISCOVER=true
            shift
            ;;
        -s|--start)
            START_SERVICES=true
            shift
            ;;
        -f|--force)
            FORCE_OVERWRITE=true
            shift
            ;;
        -b|--backup)
            CREATE_BACKUP=true
            shift
            ;;
        *)
            echo "Unknown option: $1" >&2
            show_help
            exit 1
            ;;
    esac
done

# Get script directory in a platform-independent way
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Check if required scripts exist
CONSOLIDATE_SCRIPT="$SCRIPT_DIR/consolidate-env-files.sh"
UPDATE_SCRIPT="$SCRIPT_DIR/update-docker-compose.sh"
CONFIG_YAML="$SCRIPT_DIR/services-config.yaml"

if [ ! -f "$CONSOLIDATE_SCRIPT" ]; then
    echo "Error: Consolidation script not found at $CONSOLIDATE_SCRIPT"
    exit 1
fi

if [ ! -f "$UPDATE_SCRIPT" ]; then
    echo "Error: Update script not found at $UPDATE_SCRIPT"
    exit 1
fi

# Check if the default config file exists (if no custom config provided)
if [ -z "$CONFIG_FILE" ] && [ "$AUTO_DISCOVER" = false ] && [ ! -f "$CONFIG_YAML" ]; then
    echo "Error: Default configuration file not found at $CONFIG_YAML"
    echo "Please create a configuration file or use the -d option to auto-discover services"
    exit 1
fi

# Make scripts executable if they aren't already
chmod +x "$CONSOLIDATE_SCRIPT"
chmod +x "$UPDATE_SCRIPT"

# Prepare common options for scripts
COMMON_OPTS=""
if [ "$FORCE_OVERWRITE" = true ]; then
    COMMON_OPTS="$COMMON_OPTS -f"
fi

# Prepare consolidate script options
CONSOLIDATE_OPTS="$COMMON_OPTS"
if [ "$AUTO_DISCOVER" = true ]; then
    CONSOLIDATE_OPTS="$CONSOLIDATE_OPTS -d"
fi
if [ -n "$CONFIG_FILE" ]; then
    CONSOLIDATE_OPTS="$CONSOLIDATE_OPTS -c $CONFIG_FILE"
fi

# Prepare update script options
UPDATE_OPTS="$COMMON_OPTS"
if [ "$CREATE_BACKUP" = true ]; then
    UPDATE_OPTS="$UPDATE_OPTS -b"
fi

# Step 1: Consolidate environment files
echo "Step 1: Consolidating service-specific .env files..."
"$CONSOLIDATE_SCRIPT" $CONSOLIDATE_OPTS
if [ $? -ne 0 ]; then
    echo "Error: Failed to consolidate environment files"
    exit 1
fi

# Step 2: Update docker-compose.yml
echo "Step 2: Updating docker-compose.yml with consolidated environment variables..."
"$UPDATE_SCRIPT" $UPDATE_OPTS
if [ $? -ne 0 ]; then
    echo "Error: Failed to update docker-compose.yml"
    exit 1
fi

# Step 3: Start services if requested
if [ "$START_SERVICES" = true ]; then
    echo "Step 3: Starting services with docker-compose..."

    # Check if docker-compose is installed
    if ! command -v docker-compose &> /dev/null; then
        echo "Error: docker-compose is not installed or not in the PATH"
        exit 1
    fi

    # Detect platform and use the appropriate command
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        echo "Detected macOS platform"
        docker-compose -f "$SCRIPT_DIR/docker-compose.yml" up -d
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        echo "Detected Linux platform"
        docker-compose -f "$SCRIPT_DIR/docker-compose.yml" up -d
    else
        # Other platform, try generic approach
        echo "Unknown platform, attempting to run docker-compose"
        docker-compose -f "$SCRIPT_DIR/docker-compose.yml" up -d
    fi

    if [ $? -ne 0 ]; then
        echo "Error: Failed to start services with docker-compose"
        exit 1
    fi

    echo "Services started successfully"
else
    echo "Services not started. Use -s or --start option to start services."
fi

echo "Environment workflow completed successfully!"