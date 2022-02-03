.EXPORT_ALL_VARIABLES:

BIN_FILE="/opt/calendar/calendar-app"
MIGRATION_DIR="/opt/calendar/migrations/"

DB_USER=hwuser
DB_PASSWORD=hwpasswd

OUTER_CALENDAR_SERVER_PORT = 8080
OUTER_CALENDAR_GRPC_PORT = 9090
CALENDAR_SERVER_HOST=0.0.0.0
CALENDAR_SERVER_PORT=8082
CALENDAR_GRPC_HOST=0.0.0.0
CALENDAR_GRPC_PORT=9092
CALENDAR_LOGGER_LEVEL=DEBUG

RABBITMQ_DEFAULT_USER=guest
RABBITMQ_DEFAULT_PASS=guest

generate:
	go generate ./...

gen-go: tests/integration/generated/openapicli/api_default.go

tests/integration/generated/openapicli/api_default.go: api/calendar.yaml
	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli:v5.4.0 generate --package-name=openapicli --global-property apis,models,supportingFiles=client.go -i /local/api/calendar.yaml -g go -o /local/tests/integration/generated/openapicli
	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli:v5.4.0 generate --package-name=openapicli --global-property supportingFiles=configuration.go -i /local/api/calendar.yaml -g go -o /local/tests/integration/generated/openapicli

BIN := "./bin/calendar"
BIN_SCH := "./bin/calendar_scheduler"
BIN_SEND := "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCH) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	go build -v -o $(BIN_SEND) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile_calendar .

run-img: build-img
	docker run $(DOCKER_IMG)

compose-build:
	docker-compose -f deployments/docker-compose.yaml -p calendar build

up:
	docker-compose -f deployments/docker-compose.yaml -p calendar up

down:
	docker-compose -f deployments/docker-compose.yaml -p calendar down --rmi local

version: build
	$(BIN) version

test:
	go test -race ./internal/...
	go test -race ./plugins/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

integration-tests:
	./scripts/integration_tests.sh

.PHONY: build run build-img run-img version test lint
