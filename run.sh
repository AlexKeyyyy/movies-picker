#!/usr/bin/env bash
set -euo pipefail

# 1) Остановить и удалить все контейнеры + тома
docker compose down -v

# 2) Поднять только БД
docker compose up -d postgres

echo "Waiting for Postgres to be healthy…"
# ждём, пока контейнер postgres не станет healthy
PG=$(docker compose ps -q postgres)
until [ "$(docker inspect --format='{{.State.Health.Status}}' "$PG")" = "healthy" ]; do
  echo "Postgres status: $(docker inspect --format='{{.State.Health.Status}}' "$PG")"
  sleep 2
done

# 3) Собираем образы
docker compose build builder api frontend

# 4) Запуск тестов
echo "=== Running all tests ==="
docker compose run --rm builder

# 5) Запустить API и фронтенд
echo "=== Starting API and Frontend ==="
docker compose up -d api frontend

echo
echo "✅ Done!"
echo "API:      http://localhost:${PORT:-8080}"
echo "Frontend: http://localhost:3000"
