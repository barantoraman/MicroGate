# ==================================================================================== #
# HELPERS                                        
# ==================================================================================== #
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT 								   
# ==================================================================================== #

db/migrate/new/auth:
	migrate create -seq -ext=.sql -dir=./internal/auth/db/migrations/ create_user_table

db/migrate/up/auth: confirm
	migrate -path=./internal/auth/db/migrations -database=${AUTH_DB_DSN} up

proto/create/task: confirm
	cd internal/task/pb; ./proto.sh

proto/create/auth: confirm
	cd internal/auth/pb; ./proto.sh

# Run tests for the Task service
test/task:
	go test -v -cover ./internal/task/db/mock ./internal/task/repo/mock

# Run tests for the Auth service
test/auth:
	go test -v -cover ./internal/auth/db/mock ./internal/auth/cache/mock ./internal/auth/repo/mock

# Run all service tests (Auth + Task)
test/all:
	CGO_ENABLED=0 go test -v -cover ./internal/auth/db/mock ./internal/auth/cache/mock ./internal/auth/repo/mock && \
	CGO_ENABLED=0 go test -v -cover ./internal/task/db/mock ./internal/task/repo/mock

# ==================================================================================== #
# PRODUCTION								           
# ==================================================================================== #

# Build containers for production
docker/prod/build:
	docker compose -f deployments/compose/docker-compose.yaml build

# Run containers for production
docker/prod/run:
	docker compose -f deployments/compose/docker-compose.yaml up -d

# Stop production containers
docker/prod/stop:
	docker compose -f deployments/compose/docker-compose.yaml down

# Build containers for integration tests
docker/test/build:
	docker compose -f deployments/compose/docker-compose-test.yaml build

# Run containers for integration tests
docker/test/run:
	docker compose -f deployments/compose/docker-compose-test.yaml up --build --exit-code-from integration_tests

# Stop integration test containers
docker/test/stop:
	docker compose -f deployments/compose/docker-compose-test.yaml down

# Enable Docker-related configuration in the service file
config/docker:
	sed -i.bak '/use with docker/s/^#//g;/use without docker/s/^/#/g' service.dev.yaml

# Disable Docker-related configuration for local development
config/local:
	sed -i.bak '/use with docker/s/^/#/g;/use without docker/s/^#//g' service.dev.yaml

# Start the task service locally
local/run/task:
	go run ./cmd/task

# Start the authentication service locally
local/run/auth:
	go run ./cmd/auth

# Start the API gateway service locally
local/run/gateway:
	go run ./cmd/gateway

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy -compat=1.17
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor
