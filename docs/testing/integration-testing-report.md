# Отчёт по интеграционному тестированию (Весенний этап)

## 1. Цель интеграционного тестирования

Цель — проверить совместную работу модулей backend-системы в реальном окружении приложения:
- auth + middleware + users,
- movies + repository,
- watchlist + ratings + auth,
- взаимодействие API по HTTP в поднятом окружении.

Интеграционные тесты выполняются против запущенного API и БД PostgreSQL.

## 2. Соответствие требованиям преподавателя

### 2.1 Предварительный тест-план со сценариями
В проекте реализованы сценарии в тестах `test/integration/*`.
Ниже приведён полный список реализованных интеграционных кейсов (14+):

1. Register пользователя в `TestMain`.
2. Login пользователя в `TestMain`.
3. `GET /users/me` с токеном → 200.
4. `GET /users/me` без токена → 401 (негативный).
5. `GET /movies?page=1&size=5` → 200 + проверка пагинации.
6. `GET /movies/{id}` по валидному ID → 200.
7. `GET /movies/not-a-number` → 400 (негативный).
8. `GET /movies/search` без `q` → 400 (негативный).
9. `POST /users/{id}/watchlist` валидный payload → 201.
10. `GET /users/{id}/watchlist` содержит добавленный фильм.
11. `DELETE /users/{id}/watchlist/{movieID}` → 204.
12. `POST /users/{id}/watchlist` invalid payload → 400 (негативный).
13. `POST /users/{id}/ratings` валидный payload → 201.
14. `GET /users/{id}/ratings` содержит установленный рейтинг.
15. `DELETE /users/{id}/ratings/{movieID}` → 204.
16. `POST /users/{id}/ratings` invalid payload → 400 (негативный).

**Статус:** ✅ минимум 10 сценариев выполнен с запасом.

### 2.2 Mock/stub для изоляции окружения
- В текущей реализации интеграционные тесты запускаются в docker-compose окружении с реальными API/DB сервисами.
- Внешние API не стабаются как отдельные mock-сервисы на уровне integration suite.

**Статус:** ⚠ частично.

Рекомендация для полного соответствия формулировке:
- подключить mock-сервер (например WireMock/httptest reverse proxy) для Kinopoisk/Youtube в integration pipeline,
- или зафиксировать режим запуска без внешнего API (предзагрузка seed-данных + отключение импортера).

### 2.3 Негативные сценарии
Негативные сценарии покрыты:
- unauthorized,
- invalid payload,
- invalid path/query params.

**Статус:** ✅ выполнено.

### 2.4 Запуск на CI по событию изменения кода
Запуск интеграционного пайплайна настроен в `.github/workflows/test.yml` на `push` и `pull_request`.

**Статус:** ✅ выполнено.

## 3. Использованные инструменты

- Go test + testify
- Docker Compose
- PostgreSQL
- GitHub Actions

## 4. Реализация в проекте

Основные файлы:
- `test/integration/setup_test.go` — bootstrap тестового пользователя (register/login) и подготовка token/userID.
- `test/integration/helpers_test.go` — helper для авторизованных HTTP-запросов.
- `test/integration/auth_flow_test.go` — auth/middleware сценарии.
- `test/integration/movie_flow_test.go` — сценарии movies/list/search.
- `test/integration/watchlist_flow_test.go` — lifecycle watchlist.
- `test/integration/ratings_flow_test.go` — lifecycle ratings.

CI workflow:
- `.github/workflows/test.yml`, job `integration-tests`.

## 5. Применённые техники тест-дизайна

1. **Use-case based scenarios** (пользовательские потоки).
2. **Классы эквивалентности**: корректные/некорректные payload.
3. **Граничные/валидационные проверки**: отсутствие query, нечисловой ID.
4. **Негативные проверки безопасности**: отсутствие токена.

## 6. Результаты прогонов

Локально интеграционные тесты требуют поднятого API/DB окружения (docker-compose).
Проверка в CI автоматизирована workflow `test.yml`.

Важно: ранее падения в CI были связаны не с тестами, а с шагом readiness-check БД (проверка `localhost:5432` при отсутствии проброса порта). Проверка исправлена на container health-check.

## 7. Процедура расширения набора интеграционных тестов

Пример: добавляется модуль `recommendations`.

1. Добавить endpoint и маршрутизацию.
2. Добавить integration flow:
   - `GET /users/{id}/recommendations` с токеном,
   - 401 без токена,
   - 500 при ошибке зависимости.
3. Добавить тестовые данные (seed) для устойчивого результата.
4. Обновить CI job (если нужны доп. сервисы/mock-серверы).
5. Актуализировать тест-план в отчёте.

## 8. Раздел для вклада участников команды

### Участник 1 (заполнено)
- ФИО:
- Реализованные integration-сценарии:
- Файлы:

### Участник 2 (заполнить)
- ФИО:
- Реализованные integration-сценарии:
- Файлы:

### Участник 3 (заполнить)
- ФИО:
- Реализованные integration-сценарии:
- Файлы:

