# Multi-stage build for Archivas v0.8.0

FROM golang:1.22-alpine AS build

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /src

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build all binaries (optimized)
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/archivas-node ./cmd/archivas-node && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/archivas-timelord ./cmd/archivas-timelord && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/archivas-farmer ./cmd/archivas-farmer && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/archivas ./cmd/archivas && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/archivas-registry ./cmd/archivas-registry && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/archivas-explorer ./cmd/archivas-explorer

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binaries from build stage
COPY --from=build /out/* /usr/local/bin/

# Copy genesis file
COPY genesis/devnet.genesis.json /app/genesis.json

# Create volumes
VOLUME ["/app/data", "/app/plots", "/app/logs", "/app/snapshots"]

# Expose ports (RPC, P2P, Registry, Explorer, Prometheus)
EXPOSE 8080 9090 8088 8082 9091

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -q --spider http://localhost:8080/healthz || exit 1

# Default: run node
ENTRYPOINT ["archivas-node"]
CMD ["--rpc", ":8080", "--p2p", ":9090", "--db", "/app/data", "--genesis", "/app/genesis.json", "--network-id", "archivas-devnet-v3", "--enable-gossip"]

