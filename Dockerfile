# Use the official Go image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download and install dependencies
RUN go mod download

# Copy the rest of the application code to the working directory
COPY . .

# Use the official Golang image as base
FROM golang:1.21-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Install gcc and other build essentials
RUN apk add --no-cache gcc musl-dev

# Copy the Go modules files
COPY go.mod go.sum ./

# Download and install Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o app ./cmd/api

# Create a new stage for the production image
FROM alpine:latest AS production

# Set the working directory inside the container
WORKDIR /app

# Copy the built executable from the build stage
COPY --from=build /app/app .

# Expose the port the application listens on
EXPOSE 8080

# Command to run the application
CMD ["./app"]
