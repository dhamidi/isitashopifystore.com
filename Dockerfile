FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go files and build script
COPY go.mod go.sum ./
COPY *.go ./
COPY tools ./tools

# Make build script executable
RUN chmod +x ./tools/build

# Build the application
RUN ./tools/build

# Create final image
FROM alpine:latest

WORKDIR /app

# Copy binary and assets from builder
COPY --from=builder /app/isitashopifystore .
COPY html/ ./html/
COPY assets/ ./assets/

# Expose port
EXPOSE 8080

# Run the application
CMD ["./isitashopifystore"] 