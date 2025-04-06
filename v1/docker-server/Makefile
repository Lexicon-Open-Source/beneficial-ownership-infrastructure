# Makefile for Docker Compose management
.PHONY: up down restart logs clean ps network nginx-reload nginx-test use-nginx use-traefik help

# Default Docker Compose file
COMPOSE_FILE ?= docker-compose-nginx.yml
TRAEFIK_FILE = docker-compose.yml

# Help command
help:
	@echo "Available commands:"
	@echo "  make up              - Start all services with current config"
	@echo "  make down            - Stop all services"
	@echo "  make restart         - Restart all services"
	@echo "  make logs            - View logs from all containers"
	@echo "  make logs-nginx      - View Nginx logs only (nginx mode)"
	@echo "  make logs-bo-chat    - View bo-chat logs only"
	@echo "  make logs-traefik    - View Traefik logs only (traefik mode)"
	@echo "  make ps              - List running containers"
	@echo "  make network         - Show network info"
	@echo "  make clean           - Remove containers, networks, and volumes"
	@echo "  make nginx-reload    - Reload Nginx config without downtime"
	@echo "  make nginx-test      - Test Nginx configuration"
	@echo "  make use-nginx       - Switch to Nginx configuration"
	@echo "  make use-traefik     - Switch to Traefik configuration"
	@echo "  make traefik-up      - Start Traefik services directly"
	@echo "  make nginx-up        - Start Nginx services directly"

# Start services
up:
	@echo "Starting services with $(COMPOSE_FILE)..."
	docker compose -f $(COMPOSE_FILE) up -d

# Stop services
down:
	@echo "Stopping services..."
	docker compose -f $(COMPOSE_FILE) down

# Restart services
restart:
	@echo "Restarting services..."
	docker compose -f $(COMPOSE_FILE) restart

# View logs
logs:
	docker compose -f $(COMPOSE_FILE) logs -f

# View Nginx logs only
logs-nginx:
	docker compose -f docker-compose-nginx.yml logs -f nginx

# View bo-chat logs only
logs-bo-chat:
	docker compose -f $(COMPOSE_FILE) logs -f bo-chat

# View Traefik logs only
logs-traefik:
	docker compose -f $(TRAEFIK_FILE) logs -f traefik

# List running containers
ps:
	docker compose -f $(COMPOSE_FILE) ps

# Show network info
network:
	docker network inspect traefik-public

# Clean up
clean:
	@echo "Removing containers, networks, and volumes..."
	docker compose -f $(COMPOSE_FILE) down --volumes --remove-orphans

# Reload Nginx config
nginx-reload:
	@echo "Reloading Nginx configuration..."
	docker compose -f docker-compose-nginx.yml exec nginx nginx -s reload

# Test Nginx config
nginx-test:
	@echo "Testing Nginx configuration..."
	docker compose -f docker-compose-nginx.yml exec nginx nginx -t

# Switch to Nginx configuration
use-nginx:
	@echo "Switching to Nginx configuration..."
	@echo "COMPOSE_FILE = docker-compose-nginx.yml" > .compose-file
	@echo "Now using Nginx configuration. Run 'make up' to start services."
	@export COMPOSE_FILE=docker-compose-nginx.yml

# Switch to Traefik configuration
use-traefik:
	@echo "Switching to Traefik configuration..."
	@echo "COMPOSE_FILE = docker-compose.yml" > .compose-file
	@echo "Now using Traefik configuration. Run 'make up' to start services."
	@export COMPOSE_FILE=docker-compose.yml

# Start Traefik services directly
traefik-up:
	@echo "Starting Traefik services..."
	docker compose -f $(TRAEFIK_FILE) up -d

# Start Nginx services directly
nginx-up:
	@echo "Starting Nginx services..."
	docker compose -f docker-compose-nginx.yml up -d