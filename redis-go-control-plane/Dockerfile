FROM golang:1.21 as golang-grpc-server-builder
WORKDIR /build/envoy
# Resolve and build Go dependencies as Docker cache
COPY go.mod /build/envoy/go.mod
COPY go.sum /build/envoy/go.sum
ENV GO111MODULE=on
RUN go mod download

RUN mkdir -p /build/envoy/redis-go-control-plane/example
COPY ./redis-go-control-plane/example/*.go /build/envoy/redis-go-control-plane/example/
COPY ./redis-go-control-plane/main.go /build/envoy/main.go
RUN ls -la /build/envoy/redis-go-control-plane
# Build for linux
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN go build -o server /build/envoy/.

FROM debian:bullseye-slim@sha256:9bec46ecd98ce4bf8305840b021dda9b3e1f8494a0768c407e2b233180fa1466 as os-base
RUN rm -f /etc/apt/apt.conf.d/docker-clean \
    && echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' | tee /etc/apt/apt.conf.d/keep-cache

FROM os-base as golang-grpc-server
WORKDIR /root/
# Copy the linux amd64 binary
COPY --from=golang-grpc-server-builder /build/envoy/server /bin/
ENTRYPOINT /bin/server
