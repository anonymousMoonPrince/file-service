LOCAL_BIN=$(CURDIR)/bin
PROJECT_NAME=file-service

GOLANGCI_BIN=$(LOCAL_BIN)/golangci-lint
$(GOLANGCI_BIN):
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint

GOOSE_BIN=$(LOCAL_BIN)/goose
GOOSE_TAG=v2.6.0
GOOSE_URL=https://github.com/pressly/goose/releases/download/$(GOOSE_TAG)/goose-$(shell uname -s | tr '[:upper:]' '[:lower:]')64
$(GOOSE_BIN):
	mkdir -p bin
	wget -O $(GOOSE_BIN) $(GOOSE_URL)
	chmod +x $(GOOSE_BIN)

.PHONY: up
up:
	docker-compose up -d

.PHONY: build
build:
	$(GOENV) CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -v -ldflags "$(LDFLAGS)" -o $(LOCAL_BIN)/$(PROJECT_NAME) ./cmd

.PHONY: run
run:
	$(GOENV) go run cmd/main.go $(RUN_ARGS)

.PHONY: run-checker
run-checker:
	$(GOENV) go run internal/tools/checker/main.go $(RUN_ARGS)

.PHONY: lint
lint: $(GOLANGCI_BIN)
	$(GOENV) $(GOLANGCI_BIN) run ./...

.PHONY: test
test:
	$(GOENV) go test -race -short ./...

# Migration commands
MIGRATIONS_DIR := "migrations"
POSTGRES_DSN := "postgresql://localhost:5432/file-service?user=root&password=root&sslmode=disable" # Only for local development. Delete in the future

.PHONY: migrate-up
migrate-up: $(GOOSE_BIN)
	$(GOOSE_BIN) -dir $(MIGRATIONS_DIR) postgres $(POSTGRES_DSN) up

.PHONY: migrate-down
migrate-down: $(GOOSE_BIN)
	$(GOOSE_BIN) -dir $(MIGRATIONS_DIR) postgres $(POSTGRES_DSN) down

.PHONY: migrate-reset
migrate-reset: $(GOOSE_BIN)
	$(GOOSE_BIN) -dir $(MIGRATIONS_DIR) postgres $(POSTGRES_DSN) reset

.PHONY: migrate-generate
migrate-generate: $(GOOSE_BIN)
	$(GOOSE_BIN) -dir $(MIGRATIONS_DIR) create $(name) sql

.PHONY: migrate-status
migrate-status: $(GOOSE_BIN)
	$(GOOSE_BIN) -dir $(MIGRATIONS_DIR) postgres $(POSTGRES_DSN) status
