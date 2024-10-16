.PHONY: build
build:
	@echo "Building..."
	@go run ./build -task build

.PHONY: test
test:
	@echo "Running tests..."
	@go run ./build -task test

.PHONY: publish
publish:
	@echo "Building container..."
	@go run ./build
