# Use the official Golang image as a base for building the application
FROM golang:1.22-alpine AS builder

# Install necessary packages for building
RUN apk add --no-cache git ca-certificates tzdata

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies efficiently
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .


# Build the application
# CGO_ENABLED=0 is important for creating static binaries without external dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /main ./cmd/server


# Use a minimal base image for the final, smaller runtime image
FROM alpine:latest

# Install ca-certificates for secure HTTPS connections
RUN apk --no-cache add ca-certificates wget

# Create app user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set the current working directory in the final image
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /main ./main

# Copy the OpenAPI specification file
COPY --from=builder /app/api ./api



# Change ownership to app user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port (Cloud Run uses 8080, local dev uses 8787)
EXPOSE 8080

# Health check (uses PORT env var, defaults to 8080 for Cloud Run)
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider --timeout=10 http://localhost:${PORT:-8080}/ || exit 1

# Command to run the application when the container starts
CMD ["./main"]