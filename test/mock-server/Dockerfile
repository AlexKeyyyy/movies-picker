FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o mock-server ./test/mock-server/main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/mock-server /app/mock-server
EXPOSE 9090
ENTRYPOINT ["/app/mock-server"]