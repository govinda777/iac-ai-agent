# Build stage
FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git make

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o iac-ai-agent ./cmd/agent

# Runtime stage
FROM python:3.11-alpine

# Install runtime dependencies
RUN apk add --no-cache ca-certificates curl terraform

# Install checkov
RUN pip install --no-cache-dir checkov

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/iac-ai-agent .

# Copy config
COPY configs/app.yaml ./configs/

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

ENTRYPOINT ["./iac-ai-agent"]
