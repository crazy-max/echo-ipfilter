# syntax=docker/dockerfile:1

ARG GO_VERSION="1.23"

FROM golang:${GO_VERSION}-alpine AS base
ENV GOFLAGS="-buildvcs=false"
RUN apk add --no-cache gcc linux-headers musl-dev
WORKDIR /src

FROM base AS test
ARG GO_VERSION
ENV CGO_ENABLED=1
RUN --mount=type=bind,target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build <<EOT
  set -ex
  go test -v -coverprofile=/tmp/coverage.txt -covermode=atomic -race ./...
  go tool cover -func=/tmp/coverage.txt
EOT

FROM scratch AS test-coverage
COPY --from=test /tmp/coverage.txt /
