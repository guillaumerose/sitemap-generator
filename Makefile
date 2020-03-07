all: lint test build

clean:
	rm -rf bin/

bin:
	mkdir -p bin/

.PHONY: bin/sitemap-generator
bin/sitemap-generator: bin
	go build -o bin/sitemap-generator ./cmd/standalone

.PHONY: bin/client
bin/client: bin
	go build -o bin/client ./cmd/client

.PHONY: bin/server
bin/server: bin
	go build -o bin/server ./cmd/server

.PHONY: build
build: bin/sitemap-generator bin/client bin/server

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test ./...

.PHONY: run
run: build
	bin/sitemap-generator https://www.redhat.com

images:
	docker build -t guillaumerose/sitemap-generator-server:v1.0 -f cmd/server/Dockerfile .
	docker build -t guillaumerose/sitemap-generator-client:v1.0 -f cmd/client/Dockerfile .
