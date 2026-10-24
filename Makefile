.PHONY: help help-all setup build run test types lint check

.DEFAULT_GOAL := help

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "  check   - run lint + types (tests not yet written)"
	@echo "  build   - build binary to target/"
	@echo "  run     - build and run (pass args via ARGS=\"...\")"
	@echo ""
	@echo "Run 'make help-all' for the full list, including internal targets."

help-all:
	@echo "Public targets:"
	@echo "  check   - run lint + types (tests not yet written)"
	@echo "  build   - build binary to target/"
	@echo "  run     - build and run (pass args via ARGS=\"...\")"
	@echo ""
	@echo "Internal targets (typically called via check):"
	@echo "  lint    - run linter"
	@echo "  types   - run type checker (go vet)"
	@echo "  setup   - go mod download"

setup:
	go mod download

build:
	mkdir -p target
	go build -o target/consul-ip-finder cmd/main.go

run: build
	./target/consul-ip-finder $(or $(ARGS),--help)

types:
	go vet ./...

lint:
	golangci-lint run ./...

check: lint types
	@echo "ALL CHECKS PASSED"
