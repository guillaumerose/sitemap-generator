# sitemap-generator

## Build

Please install go and run the following command:

```
make build
```

This will output the `sitemap-generator` binary in `bin/`

## Usage

Standalone crawler:

```
$ ./bin/sitemap-generator https://www.redhat.com
/en
/en/search
/en/technologies
/en/solutions
...
```

client/server:

```
$ ./bin/server
INFO[0000] Listening on :8080
```

```
$ ./bin/client https://kompose.io
INFO[0000] Crawling https://kompose.io (parallelism: 10, maxDepth: 5)
INFO[0000] 7 URLs found
- /
- /architecture
- /conversion
- /docs
  - /conversion.md
  - /maven-example.md
- /getting-started
- /installation
- /integrations
- /user-guide
```
