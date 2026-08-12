# Makefile for a Gin-based Golang project using gox for cross-compilation

# Project Name
BINARY_NAME=fernclient

# Go related variables.
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOPKG=$(GOBASE)

mod-tidy:
	@echo "🧹 Running go mod tidy..."
	@go mod tidy

unit-test:
	@echo "🧪 Running Unit Tests..."
	ginkgo -r -p --label-filter=unit --randomize-all

test:
	@echo "🧪 Running All Tests with labels \"$(LABEL_FILTER)\"..."
	ginkgo -r -p --label-filter="$(LABEL_FILTER)" --randomize-all

# Requires a Fern Platform instance you can reach directly (e.g. running
# locally via `make docker-test-up` in a fern-platform checkout). There's no
# hosted or CI-reachable Fern instance for this open-source project yet, so
# this target only works against your own local instance for now.
report-local: ## Run the dogfood suite and report results to a local Fern (expects FERN_BASE_URL and PROJECT_ID set)
	@echo "📡 Reporting to $${FERN_BASE_URL:-http://localhost:8080/} as project $${PROJECT_ID}..."
	ginkgo -r ./tests

fmt:
	@echo "📝 Formatting Go code..."
	@go fmt ./...

lint: 
	@echo "🔍 Running linter..."
	@golangci-lint run ./...