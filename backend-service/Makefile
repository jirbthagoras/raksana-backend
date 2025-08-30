DOCKER_COMPOSE = docker compose

# Makefile increases dev experience for continous development!

# Default target
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make build        - Build the Docker Compose project"
	@echo "  make up           - Run containers in detached mode"
	@echo "  make down         - Stop and remove containers"
	@echo "  make restart      - Restart containers"
	@echo "  make logs         - Show logs"
	@echo "  make prune        - Remove dangling images"
	@echo "  make clean        - Stop all and remove images/volumes"

.PHONY: build
build:
	$(DOCKER_COMPOSE) build --no-cache

.PHONY: up
up:
	$(DOCKER_COMPOSE) up -d

.PHONY: down
down:
	@$(DOCKER_COMPOSE) down --remove-orphans

.PHONY: restart
restart: down up

.PHONY: logs
logs:
	$(DOCKER_COMPOSE) logs -f

.PHONY: prune
prune:
	docker image prune -f

.PHONY: clean
clean:
	@dangling=$$(docker images -f "dangling=true" -q); \
	if [ -n "$$dangling" ]; then \
		docker rmi $$dangling; \
	else \
		echo "No dangling images to remove."; \
	fi
