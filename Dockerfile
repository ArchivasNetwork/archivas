# Multi-stage build for Archivas

FROM golang:1.22-alpine AS build

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /src

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build all binaries
RUN go build -o /out/archivas-node ./cmd/archivas-node && \
    go build -o /out/archivas-timelord ./cmd/archivas-timelord && \
    go build -o /out/archivas-farmer ./cmd/archivas-farmer && \
    go build -o /out/archivas-wallet ./cmd/archivas-wallet

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binaries from build stage
COPY --from=build /out/* /usr/local/bin/

# Copy genesis file
COPY genesis/devnet.genesis.json /app/genesis.json

# Create volumes
VOLUME ["/app/data", "/app/plots", "/app/logs"]

# Expose ports
EXPOSE 8080 9090

# Default command: run node
ENTRYPOINT ["archivas-node"]
CMD ["--rpc", ":8080", "--p2p", ":9090", "--db", "/app/data", "--genesis", "/app/genesis.json", "--network-id", "archivas-devnet-v3"]

