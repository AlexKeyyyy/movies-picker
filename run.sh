#!/usr/bin/env bash
set -euo pipefail

# 1) Сброс всего – убираем старые контейнеры и тома (по желанию)
docker-compose down -v

# 2) Поднимаем только БД
docker-compose up -d postgres

echo "Waiting for Postgres to be healthy…"
# ждём healthcheck
until docker-compose exec -T postgres pg_isready -U postgres; do
  sleep 2
done

# 3) Собираем образы
docker-compose build api builder

# 4) Запуск unit-тестов
echo "→ Running unit tests…"
docker-compose run --rm builder go test ./... -cover

# 5) Запуск интеграционных тестов
echo "→ Running integration tests…"
# здесь жёсткий timeout, если нужно:
docker-compose run --rm builder go test ./internal/handlers -timeout 120s -cover

# 6) Поднять API (оно само импортит фильмы перед стартом)
echo "→ Starting API with import_all_movies"
docker-compose up -d api

echo "✅ All done! Visit: http://localhost:${PORT:-8080}"
