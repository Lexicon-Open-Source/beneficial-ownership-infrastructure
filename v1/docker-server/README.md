# Traefik Reverse Proxy with Let's Encrypt SSL

This directory contains the configuration for a Traefik reverse proxy with automatic SSL certificate management via Let's Encrypt.

## Prerequisites

- Docker and Docker Compose installed
- A publicly accessible domain name pointing to your server
- Ports 80 and 443 open on your firewall

## Configuration

Before deploying:

1. Copy the environment template to create your configuration:
   ```bash
   cp .env.example .env
   ```

2. Configure the environment variables in the `.env` file:
   - `TRAEFIK_DOMAIN`: Your domain for accessing the Traefik dashboard
   - `HTTP_PORT` and `HTTPS_PORT`: Ports for HTTP and HTTPS traffic
   - `TRAEFIK_ADMIN_USER` and `TRAEFIK_ADMIN_PASSWORD`: Dashboard credentials
   - `ACME_EMAIL`: Your email for Let's Encrypt notifications
   - `TRAEFIK_LOG_LEVEL`: Log level (DEBUG, INFO, WARN, ERROR)
   - `DOCKER_NETWORK`: Name of the Docker network for Traefik

3. Set up authentication:
   ```bash
   # Generate a password hash and update configs
   ./generate-password.sh <username> <password>
   ```
   This will update both the `.env` file and the `traefik/config/users` file with your credentials.

## Deployment

```bash
# Create the Docker network (if not already created)
docker network create ${DOCKER_NETWORK}

# Start the Traefik service
docker-compose up -d
```

## Adding Services

To add a service behind the Traefik proxy, add the following labels to your service in docker-compose.yml:

```yaml
labels:
  - "traefik.enable=true"
  - "traefik.http.routers.myservice.rule=Host(`myservice.yourdomain.com`)"
  - "traefik.http.routers.myservice.entrypoints=websecure"
  - "traefik.http.routers.myservice.tls.certresolver=letsencrypt"
  - "traefik.http.services.myservice.loadbalancer.server.port=8080"  # Your service's exposed port
  # Optional: Apply secure headers and other middlewares
  - "traefik.http.routers.myservice.middlewares=secure-headers@file"
```

## Security

The configuration includes:
- Automatic HTTP to HTTPS redirection
- Let's Encrypt SSL certificates with auto-renewal
- Secure headers middleware
- Rate limiting
- Basic authentication for the Traefik dashboard

## Access Traefik Dashboard

Once deployed, access the Traefik dashboard at: https://${TRAEFIK_DOMAIN}
(Login with the credentials configured in the `traefik/config/users` file)

## Version Control

- The `.env.example` file is committed to version control as a template.
- The actual `.env` file and sensitive files like `acme.json` and the `users` file are excluded from version control in the `.gitignore`.