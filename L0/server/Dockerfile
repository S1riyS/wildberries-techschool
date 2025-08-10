# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build with optimizations
RUN CGO_ENABLED=0 GOOS=linux \
    go build -ldflags="-w -s" -trimpath \
    -o /build/server ./cmd/server/main.go

# Runtime stage
FROM alpine:3.22
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /build/server .

# Non-root user
RUN addgroup -S app \
    && adduser -S app -G app \
    && chown -R app:app /app
USER app

EXPOSE 8080
ENTRYPOINT ["./server"]