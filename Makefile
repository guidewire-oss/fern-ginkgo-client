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