# Stage 1: builder
FROM golang:1.23-alpine AS builder
WORKDIR /app

# скачиваем модули
COPY go.mod go.sum ./
RUN go mod download

# копируем весь код
COPY . .

# собираем два бинаря
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o import_all_movies ./cmd/tools/import_all_movies.go

# Stage 2: runtime
FROM alpine:3.18 AS runtime
RUN apk add --no-cache ca-certificates postgresql-client

WORKDIR /app
COPY --from=builder /app/server       /app/server
COPY --from=builder /app/import_all_movies /app/import_all_movies
COPY .env                             /app/.env

EXPOSE 8080
ENTRYPOINT ["/app/server"]
