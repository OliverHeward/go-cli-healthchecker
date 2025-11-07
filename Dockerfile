# Stage 1: Builder
FROM golang:1.24-alpine AS builder

# Install build dependencies (Alpine needs these)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files first (for layer caching)
COPY go.mod go.sum ./

# Download dependencies (cached if go.mod/go.sum don't change)
RUN go mod download

# Copy source code
COPY . .

# Build the binary
# CGO_ENABLED=0: Static binary (no C dependencies)
# -ldflags="-w -s": Strip debug info (smaller binary)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o healthcheck

# Stage 2: Runtime
FROM alpine:latest

# Install CA certificates (for HTTPS requests)
RUN apk --no-cache add ca-certificates

# Create non-root user for security
RUN addgroup -g 1000 healthcheck && \
    adduser -D -u 1000 -G healthcheck healthcheck

WORKDIR /home/healthcheck

# Copy ONLY the binary from builder stage
COPY --from=builder /app/healthcheck .

# Change ownership to non-root user
RUN chown healthcheck:healthcheck healthcheck

# Switch to non-root user
USER healthcheck

# Run the health check
ENTRYPOINT ["./healthcheck"]
CMD ["check"]