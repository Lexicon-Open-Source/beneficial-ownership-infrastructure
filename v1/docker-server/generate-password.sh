#!/bin/bash
# Script to generate a hashed password for Traefik basic auth and update configs

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <username> <password>"
    exit 1
fi

USERNAME=$1
PASSWORD=$2

# Generate the hash using openssl
HASH=$(openssl passwd -apr1 "$PASSWORD" | sed -e "s/^\\$apr1\\$/$USERNAME:&/")

echo "Generated hash: $HASH"

# Update the .env file
if [ -f .env ]; then
    # Update username and password in .env
    sed -i.bak "s/^TRAEFIK_ADMIN_USER=.*/TRAEFIK_ADMIN_USER=$USERNAME/" .env
    sed -i.bak "s/^TRAEFIK_ADMIN_PASSWORD=.*/TRAEFIK_ADMIN_PASSWORD=$PASSWORD/" .env

    echo "Updated .env file with new credentials"
else
    echo "Error: .env file not found"
    exit 1
fi

# Ensure the config directory exists
mkdir -p traefik/config

# Update the users file for Traefik auth
# Make sure the hash contains the username in the format username:hash
if [[ "$HASH" != "$USERNAME:"* ]]; then
    # If username is not already in the hash, ensure it's formatted correctly
    CLEAN_HASH="$USERNAME:$(echo "$HASH" | cut -d':' -f2)"
else
    # If username is already in the hash, just clean any escaped $ characters
    CLEAN_HASH=$(echo "$HASH" | sed -e s/\\$\\$/\\$/g)
fi

# Write to the users file
echo "$CLEAN_HASH" > traefik/config/users

echo "Updated traefik/config/users file with new credentials"

# Set proper permissions for the users file
chmod 600 traefik/config/users

echo ""
echo "Authentication credentials updated successfully."
echo "Remember to restart Traefik to apply changes:"
echo "docker-compose down && docker-compose up -d"