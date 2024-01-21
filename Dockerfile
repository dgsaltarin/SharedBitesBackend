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

# Build the Go application
RUN go build -o main .

# Expose the port that the application will run on
EXPOSE 8080

# Command to run the application
CMD ["./main"]