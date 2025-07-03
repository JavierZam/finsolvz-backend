# Use the official Golang image as a base for building the application
FROM golang:1.22-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies efficiently
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
# CGO_ENABLED=0 is important for creating static binaries without external dependencies
# -a -installsuffix cgo helps reduce the final image size
RUN CGO_ENABLED=0 go build -o /main ./cmd/server

# Use a minimal base image for the final, smaller runtime image
FROM alpine:latest

# Install ca-certificates for secure HTTPS connections
RUN apk --no-cache add ca-certificates

# Set the current working directory in the final image
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /main .

# Expose the port your application will listen on
EXPOSE 8080

# Command to run the application when the container starts
CMD ["./main"]