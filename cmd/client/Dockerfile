FROM golang:1.14-alpine3.11 as build-env
RUN apk add -u make

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make bin/client

FROM alpine:3.11
COPY --from=build-env /app/bin/client /client
ENTRYPOINT ["/client"]
