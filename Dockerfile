# Stage 1: build Go binary
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Stage 2: minimal runtime image
FROM scratch
WORKDIR /app
COPY --from=builder /app/server /app/server
COPY .env /app/.env
EXPOSE 8080
ENTRYPOINT ["/app/server"]
