all: lint test build

build:
	mkdir -p bin/
	go build -o bin/sitemap-generator ./cmd/standalone
	go build -o bin/server ./cmd/server

lint:
	golangci-lint run ./...

test:
	go test ./...

run: build
	bin/sitemap-generator https://www.redhat.com

