# run.ps1
Param(
    [int]$TimeoutSeconds = 300
)

# Остановить старая сборка
docker-compose down

# 1) Собрать образы и поднять БД
docker-compose up --build -d postgres

Write-Host "Waiting for PostgreSQL to be healthy..."
# Ждём healthcheck
Start-Sleep -Seconds 10

# 2) Собрать образ API
docker-compose build api

# 3) Запустить unit-тесты внутри контейнера
Write-Host "Running unit tests..."
docker-compose run --rm api powershell -Command "go test ./... -cover"

# 4) Запустить интеграционные тесты
Write-Host "Running integration tests..."
docker-compose run --rm api powershell -Command "go test ./internal/handlers -timeout ${TimeoutSeconds}s -cover"

# 5) Перезапустить сервис полностью
docker-compose down
docker-compose up -d

Write-Host "✅ All done! Service running on port $($env:PORT)"
