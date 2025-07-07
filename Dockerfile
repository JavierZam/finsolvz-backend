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

# ✅ DEBUG: Show what files are available in build context
RUN echo "=== BUILD CONTEXT FILES ===" && \
    ls -la . && \
    echo "=== API DIRECTORY ===" && \
    ls -la ./api/ || echo "No api directory found" && \
    echo "=== CHECKING OPENAPI FILE ===" && \
    cat ./api/openapi.yaml | head -10 || echo "No openapi.yaml found"

# Build the application
# CGO_ENABLED=0 is important for creating static binaries without external dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /main ./cmd/server

# ✅ VERIFY: Binary created successfully
RUN ls -la /main && file /main

# Use a minimal base image for the final, smaller runtime image
FROM alpine:latest

# Install ca-certificates for secure HTTPS connections and debugging tools
RUN apk --no-cache add ca-certificates wget curl

# Create app user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set the current working directory in the final image
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /main ./main

# ✅ EXPLICIT: Copy API files (even if not found, won't fail build)
COPY --from=builder /app/api ./api 2>/dev/null || echo "No api directory found in builder stage" /app/

# ✅ DEBUG: Show final container contents
RUN echo "=== FINAL CONTAINER FILES ===" && \
    ls -la . && \
    echo "=== API DIRECTORY ===" && \
    ls -la ./api/ || echo "No api directory in final container" && \
    echo "=== BINARY CHECK ===" && \
    ls -la ./main && \
    file ./main

# Change ownership to app user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider --timeout=10 http://localhost:8080/ || exit 1

# Command to run the application when the container starts
CMD ["./main"]