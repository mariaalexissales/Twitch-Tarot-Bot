# Build stage
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /bot ./cmd/bot

# Run stage
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /bot /bot
ENTRYPOINT ["/bot"]
