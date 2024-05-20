# Build the Docker images
build:
	docker-compose build

# Start the service in detached mode
up:
	docker-compose up -d

# Stop the service
down:
	docker-compose down

# View the logs
logs:
	docker-compose logs -f

# Restart the service
restart: down up

# Remove all stopped containers and unused images
clean:
	docker system prune -f

# Execute the tests
test:
	go test ./...
