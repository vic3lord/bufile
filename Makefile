IMAGE ?= bufile
TAG   ?= dev

# Prefer apple/container, fall back to docker. Override with CONTAINER_CLI=...
CONTAINER_CLI ?= $(shell command -v container 2>/dev/null || command -v docker 2>/dev/null)

.PHONY: build
build:
	@echo "Building..."
	@go build .

.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...

# Local image build only. Publishing to ghcr.io is done by CI.
.PHONY: container
container:
	@test -n "$(CONTAINER_CLI)" || { \
		echo "No container CLI found: install apple/container or docker,"; \
		echo "or set CONTAINER_CLI=/path/to/cli"; \
		exit 1; \
	}
	@echo "Building container with $(notdir $(CONTAINER_CLI))..."
	@$(CONTAINER_CLI) build -t $(IMAGE):$(TAG) .
