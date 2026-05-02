# ---------- BUILD STAGE ----------
FROM golang:1.22 AS builder

WORKDIR /build

# Copy go files first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy full source
COPY . .

# Build binary
RUN go build -o app ./cmd/main.go


# ---------- RUNTIME STAGE ----------
FROM debian:bookworm-slim

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/app /app/app

# Ensure executable
RUN chmod +x /app/app

EXPOSE 8080

# Run binary
CMD ["/app/app"]