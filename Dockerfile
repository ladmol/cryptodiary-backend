FROM golang:1.21 AS builder

# Set working directory
WORKDIR /app

# Copy go.mod, go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source files
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cryptodiary-backend .

# Final stage: Create a minimal image
FROM alpine:latest

# Add necessary packages
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/cryptodiary-backend .

# Expose the application port
EXPOSE 3001

# Command to run the application
CMD ["./cryptodiary-backend"]
