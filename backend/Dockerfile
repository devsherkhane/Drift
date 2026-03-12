# ---- Build Stage ----
FROM golang:1.18-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# ---- Run Stage ----
FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy only the binary from the build stage
COPY --from=builder /app/server .
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations

# Create uploads directory
RUN mkdir -p uploads/avatars

EXPOSE 8080

CMD ["./server"]
