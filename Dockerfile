# Multi-stage build for Raspberry Pi (ARM64)
FROM --platform=linux/arm64 golang:1.21-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Install build dependencies for ARM64
RUN apk add --no-cache \
    gcc \
    musl-dev \
    git \
    ca-certificates \
    && update-ca-certificates

# Set Go environment variables for better module handling
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org

# Copy the Go modules files first for better caching
COPY go.mod go.sum ./

# Download and cache Go module dependencies with verbose output
RUN go mod download -x

# Verify modules
RUN go mod verify

# Copy the rest of the application source code
COPY . .

# Build the Go application for ARM64
# CGO_ENABLED=1 is needed for some dependencies like SQLite
RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -o app ./cmd/api

# Production stage - minimal Alpine image for Raspberry Pi
FROM --platform=linux/arm64 alpine:latest AS production

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    wget \
    && update-ca-certificates

# Create a non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set the working directory
WORKDIR /app

# Copy the built executable from the build stage
COPY --from=build /app/app .

# Change ownership to the non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose the port the application listens on
EXPOSE 8080

# Add health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthcheck || exit 1

# Command to run the application
CMD ["./app"]
