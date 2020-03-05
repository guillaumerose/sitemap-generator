all: build

build:
	mkdir -p bin/
	go build -o bin/sitemap-generator .

run: build
	bin/sitemap-generator https://www.redhat.com
