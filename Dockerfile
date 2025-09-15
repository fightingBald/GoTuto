# syntax=docker/dockerfile:1

ARG GO_VERSION=1.24
FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git build-base

# Enable module caching
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy the whole repo (adjust if you want a narrower context)
COPY . .

# Build the service binary
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/product-query-svc ./backend/cmd/marketplace/product-query-svc

# Runtime
FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=builder /out/product-query-svc ./product-query-svc
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/product-query-svc"]
