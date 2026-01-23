# ---------- Build stage ----------
FROM golang:1.22-alpine AS builder

# Build a static Linux binary
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build CLI binary
RUN go build -o fcf ./cmd/fcf


# ---------- Runtime stage ----------
FROM alpine:latest

# Runtime dependencies
RUN apk add --no-cache bash

# Create non-root user with bash as default shell
RUN adduser -D -h /home/fcfuser -s /bin/bash fcfuser

# Copy binary from builder
COPY --from=builder /app/fcf /usr/local/bin/fcf
RUN chmod +x /usr/local/bin/fcf

# Switch to non-root user
USER fcfuser
WORKDIR /home/fcfuser

# Create .bashrc file before installing shell integration
RUN touch ~/.bashrc

# Install shell integration for navigation support (binary already in /usr/local/bin)
RUN /usr/local/bin/fcf install --shell-only

# Default: start interactive bash with shell integration loaded
CMD ["/bin/bash"]
