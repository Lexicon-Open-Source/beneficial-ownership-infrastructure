#!/bin/bash
# Initial setup script for Traefik reverse proxy

# Colors for better readability
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Setting up Traefik reverse proxy...${NC}"

# Check if .env file exists, create from example if not
if [ ! -f .env ]; then
    echo -e "${YELLOW}Creating .env file from template...${NC}"
    if [ -f .env.example ]; then
        cp .env.example .env
        echo -e "${GREEN}Created .env file. Please update it with your configuration.${NC}"
    else
        echo -e "${RED}Error: .env.example not found!${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}.env file already exists.${NC}"
fi

# Create traefik directory structure
echo -e "${YELLOW}Creating directory structure...${NC}"
mkdir -p traefik/config

# Create acme.json file with correct permissions
echo -e "${YELLOW}Creating acme.json for Let's Encrypt certificates...${NC}"
touch traefik/acme.json
chmod 600 traefik/acme.json
echo -e "${GREEN}Created acme.json with permissions 600${NC}"

# Ensure generate-password.sh is executable
echo -e "${YELLOW}Setting execute permissions on scripts...${NC}"
chmod +x generate-password.sh
echo -e "${GREEN}Made generate-password.sh executable${NC}"

# Create Docker network
echo -e "${YELLOW}Creating Docker network...${NC}"
NETWORK_NAME=$(grep DOCKER_NETWORK .env | cut -d= -f2)
if [ -z "$NETWORK_NAME" ]; then
    NETWORK_NAME="traefik-public"
    echo -e "${YELLOW}Using default network name: ${NETWORK_NAME}${NC}"
fi

docker network create $NETWORK_NAME || true
echo -e "${GREEN}Docker network '$NETWORK_NAME' is ready${NC}"

# Prompt for credentials
echo -e "${YELLOW}Do you want to set up authentication credentials now? [y/N]${NC}"
read -r response
if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    echo -e "${YELLOW}Please enter a username for Traefik dashboard:${NC}"
    read -r username
    echo -e "${YELLOW}Please enter a password for Traefik dashboard:${NC}"
    read -rs password

    echo -e "${YELLOW}Generating credentials...${NC}"
    ./generate-password.sh "$username" "$password"
else
    echo -e "${YELLOW}Skipping authentication setup. Run ./generate-password.sh manually later.${NC}"
fi

echo -e "\n${GREEN}Setup complete!${NC}"
echo -e "${YELLOW}Next steps:${NC}"
echo -e "1. Review and update the .env file with your domain and email"
echo -e "2. Run 'docker-compose up -d' to start Traefik"