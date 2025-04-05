# Start from golang base image
FROM golang:1.20-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Start a new stage from scratch for a smaller final image
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Install ca-certificates for HTTPS requests if needed
RUN apk --no-cache add ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
