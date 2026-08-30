REDIS_CONT_NAME ?= file-processor-redis-test
REDIS_PORT ?= 6380
DOCKER_USERNAME ?= rajsekhar05

.PHONY: help build-all push-all run-app stop-app run-redis stop-redis

# Default target when you just run 'make'
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build-all: ## Build Docker images for API and Worker
	docker build -f api/Dockerfile -t ${DOCKER_USERNAME}/file-processor-api:latest .
	docker build -f worker/Dockerfile -t ${DOCKER_USERNAME}/file-processor-worker:latest .

push-all: ## Push Docker images to Docker Hub
	docker push ${DOCKER_USERNAME}/file-processor-api:latest
	docker push ${DOCKER_USERNAME}/file-processor-worker:latest


run-app: ## Start the application stack via Docker Compose
	docker compose up -d

stop-app: ## Stop the application stack and remove volumes/orphans
	docker compose down --volumes --remove-orphans


run-redis: ## Start a standalone Redis container for testing
	docker run --name ${REDIS_CONT_NAME} -p ${REDIS_PORT}:6379 -d redis:alpine

stop-redis: ## Stop and remove the standalone Redis container
	docker stop ${REDIS_CONT_NAME} && docker rm ${REDIS_CONT_NAME}