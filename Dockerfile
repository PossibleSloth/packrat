# ---------- Stage 1: Build ----------
FROM golang:1.25-alpine AS builder

# Set the working directory
WORKDIR /app

# Ensure a portable, static-ish binary
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application (strip debug info for smaller size)
RUN go build -o myapp .

# ---------- Stage 2: Final ----------
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Install runtime dependencies you actually need
RUN apk add --no-cache ca-certificates tzdata


# Copy the binary and set ownership
COPY --from=builder /app/myapp /app/myapp

# Set the entrypoint command
ENTRYPOINT ["/app/myapp"]
