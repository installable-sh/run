.PHONY: build test lint fmt deps update-lib docker ci

# Update lib to latest v1 tag and build
build: update-lib
	go build -ldflags="-s -w" -o RUN .

# Update lib dependency to latest v1
update-lib:
	go get github.com/installable-sh/lib@v1
	go mod tidy

# Run tests
test: update-lib
	go test -v ./...

# Check formatting
fmt:
	@test -z "$$(gofmt -l .)" || (echo "Files not formatted:"; gofmt -l .; exit 1)

# Run linter
lint:
	golangci-lint run

# Run all CI checks
ci: update-lib fmt lint test

# Download all dependencies
deps:
	go mod download

# Build Docker image
docker: update-lib
	docker build -t installable-sh/run .
