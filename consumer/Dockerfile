FROM golang:1.24-alpine

WORKDIR /app

# Copy go.work and all module go.mod files first (for caching)
COPY go.work ./
COPY publisher/go.mod publisher/
COPY consumer/go.mod consumer/
COPY nf-server/go.mod nf-server/

# Sync workspace and download all dependencies
RUN go work sync
RUN go mod download

# Copy all modules into the container
COPY publisher ./publisher
COPY nf-server ./nf-server
COPY consumer ./consumer

# Build the binary for this module
RUN go build -o /consumer ./consumer

# Default command
CMD ["/consumer"]