<#
.SYNOPSIS
  Полный CI-процесс: поднять БД, собрать образы, прогнать тесты, импортировать фильмы, запустить API.
.DESCRIPTION
  Требует Docker Desktop (docker и docker-compose в PATH), PowerShell 7+.
#>

param(
    [int]$IntegrationTimeout = 120
)

# 1) Остановить и удалить все старые контейнеры + том с данными
docker-compose down -v

# 2) Поднять только БД
docker-compose up -d postgres

Write-Host "Waiting for PostgreSQL to be healthy…"
while (-not (docker-compose exec -T postgres pg_isready -U postgres)) {
    Start-Sleep -Seconds 2
}

# 3) Собрать образы
docker-compose build api builder

# 4) Запуск unit-тестов
Write-Host "→ Running unit tests…"
docker-compose run --rm builder pwsh -Command "go test ./... -cover"

# 5) Запуск интеграционных тестов
Write-Host "→ Running integration tests…"
docker-compose run --rm builder pwsh -Command "go test ./internal/handlers -timeout ${IntegrationTimeout}s -cover"

# 6) Запустить API (в нём при старте выполнится import_all_movies)
Write-Host "→ Starting API with import_all_movies…"
docker-compose up -d api

Write-Host "✅ All done! Service running on port $($env:PORT)"
