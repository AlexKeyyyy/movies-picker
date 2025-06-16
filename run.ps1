param(
  [int]$IntegrationTimeout = 120
)

# 1) Сброс
docker compose down -v

# 2) Поднять только postgres
docker compose up -d postgres

Write-Host "Waiting for Postgres to be healthy…"
# получаем ID контейнера postgres
$pg = docker compose ps -q postgres
do {
  Start-Sleep -Seconds 2
  $status = docker inspect --format='{{.State.Health.Status}}' $pg
  Write-Host "Postgres status: $status"
} while ($status -ne "healthy")

# 3) Собрать образы
docker compose build builder api frontend

# 4) Запустить тесты
Write-Host "→ Running all tests…"
docker compose run --rm builder

# 5) Запустить API и фронтенд
Write-Host "→ Starting API and Frontend…"
docker compose up -d api frontend

Write-Host "✅ Done! API: http://localhost:$($env:PORT); Frontend: http://localhost:3000"
