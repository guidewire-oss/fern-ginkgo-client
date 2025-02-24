# Makefile for a Gin-based Golang project using gox for cross-compilation

# Project Name
BINARY_NAME=fernclient

# Go related variables.
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOPKG=$(GOBASE)

# Testing

unit-test:
	@echo "🧪 Running Unit Tests..."
#	@go test $(TEST_FLAGS) -coverprofile=profile.cov ./...
	ginkgo -r -p --label-filter=unit --randomize-all

test:
	@echo "🧪 Running All Tests with labels \"$(LABEL_FILTER)\"..."
#	@go test $(TEST_FLAGS) -coverprofile=profile.cov ./...
	ginkgo -r -p --label-filter="$(LABEL_FILTER)" --randomize-all