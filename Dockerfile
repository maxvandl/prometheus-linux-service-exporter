# Use the official Golang image to build the application
FROM golang:1.24 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go binary
RUN go build -o xrdp_monitor

# Create a minimal runtime image
FROM ubuntu:24.04

# Set the working directory inside the runtime container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/xrdp_monitor .

# Expose the monitoring port
EXPOSE 8888

# Run the binary
CMD ["./xrdp_monitor"]
