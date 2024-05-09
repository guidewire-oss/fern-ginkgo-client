# Makefile for a Gin-based Golang project using gox for cross-compilation

# Project Name
BINARY_NAME=fernclient

# Go related variables.
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOPKG=$(GOBASE)

# Testing
test:
	@echo "Testing Client..."
	@go test ./...
