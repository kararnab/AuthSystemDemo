# ---- Build stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install CA certs for HTTPS (needed for go mod)
RUN apk add --no-cache ca-certificates

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

# Build a static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o auth-demo-go ./cmd/server

# ---- Runtime stage ----
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy binary only
COPY --from=builder /app/auth-demo-go /app/auth-demo-go

# Run as non-root
USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/app/auth-demo-go"]
