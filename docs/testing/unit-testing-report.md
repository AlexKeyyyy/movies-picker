# Отчёт по модульному тестированию (Весенний этап)

> Документ подготовлен для проверки преподавателем.
> Ниже зафиксированы: соответствие требованиям, инструменты, техники тест-дизайна, результаты запусков, покрытие и порядок расширения тестового набора.

## 1. Цель и область модульного тестирования

Цель модульного тестирования — проверить корректность отдельных модулей backend-части (репозиторий, middleware, API-валидация handlers, клиенты внешних API) в изоляции и зафиксировать поведение на позитивных и негативных данных.

Покрытые уровни:
- `internal/repository` — CRUD и бизнес-операции хранения данных (users, movies, watchlist, ratings).
- `internal/middleware` — JWT-проверка и извлечение user_id из токена.
- `internal/handlers` — валидация входных параметров/JSON и корректные HTTP-коды ошибок.
- `pkg/kinopoisk`, `pkg/youtube` — клиентская логика внешних API и обработка ошибок.

## 2. Требования преподавателя и соответствие

### 2.1 Автоматический запуск при сборке/в CI
- В проекте настроен GitHub Actions workflow для unit-тестов: `.github/workflows/test.yml`, job `unit-tests`.
- Unit-тесты запускаются командой `go test ./internal/... ./pkg/... -v -cover`.
- Запуск привязан к `push` и `pull_request` по всем веткам.

**Статус:** ✅ выполнено.

### 2.2 Минимум 25 unit-тестов на участника
- На текущей ветке найдено:
  - 36 тестов в `internal/*_test.go`;
  - 44 теста в совокупности `internal + pkg`.

**Статус:** ✅ порог 25+ достигнут.

### 2.3 Применение техник тест-дизайна
Применены следующие техники:
1. **Классы эквивалентности**:
   - валидный/невалидный JWT;
   - валидный/невалидный JSON payload;
   - существующий/несуществующий объект в репозитории.
2. **Граничные условия**:
   - некорректные значения path/query параметров (`movie_id`, `q`);
   - пустые/отсутствующие заголовки авторизации.
3. **Негативные сценарии**:
   - malformed token, wrong signature, expired token;
   - invalid JSON в handlers;
   - некорректные форматы идентификаторов в URL.

**Статус:** ✅ выполнено.

## 3. Использованные инструменты

- Язык/фреймворк: Go + `testing`.
- Assertions: `testify` (`assert`, `require`).
- Покрытие: `go test -coverprofile` + `go tool cover -func`.
- CI: GitHub Actions.

## 4. Структура текущих unit-тестов

### 4.1 Repository
Файл: `internal/repository/repo_test.go`

Покрытые операции:
- Users: `CreateUser`, `GetUserByEmail`, `GetUserByID`, `UpdateUser`.
- Movies: `UpsertMovie`, `GetMovieByID`, `SearchMovies`, `ListMovies`, `ListPopularMovies`.
- Watchlist: `AddToWatchlist`, `GetWatchlist`, `RemoveFromWatchlist`.
- Ratings: `UpsertRating`, `GetRatings`, `DeleteRating`.

### 4.2 Middleware
Файлы:
- `internal/middleware/auth_test.go`
- `internal/middleware/auth_additional_test.go`

Покрытые сценарии:
- валидный токен;
- отсутствующий `Authorization`;
- неверный формат `Bearer`;
- пустой bearer token;
- expired token;
- wrong signature;
- malformed token;
- некорректный тип/отсутствие `user_id` в claims;
- проверка прокидывания `user_id` в context.

### 4.3 Handlers (валидация)
Файл: `internal/handlers/handlers_validation_test.go`

Покрытые сценарии:
- invalid JSON для `auth/register`, `auth/login`, `users/me`, `watchlist`, `ratings`;
- invalid `id` для `/movies/{id}` и `/movies/{id}/reviews`;
- missing query для `/movies/search`;
- проверка конструкторов обработчиков.

### 4.4 Внешние API клиенты
Файлы:
- `pkg/kinopoisk/client_test.go`
- `pkg/youtube/client_test.go`

Покрыты позитивные и негативные ветки клиентской логики.

## 5. Результаты запусков и покрытие

Команды локальной проверки:
1. `go test ./internal/... ./pkg/... -coverprofile=coverage.out`
2. `go tool cover -func=coverage.out | tail -n 1`

Текущий результат total coverage:
- **46.9%** (в рамках `internal + pkg`).

Пакетные ориентиры:
- `internal/middleware`: 100%
- `internal/repository`: 88.9%
- `internal/handlers`: 33.6%
- `internal/service`: 0%
- `pkg/kinopoisk`: 86.5%
- `pkg/youtube`: 96.7%

### Вывод по покрытию
- Требование **80%** по всей кодовой базе на текущем этапе **не достигнуто**.
- Главный непокрытый слой: `internal/service`.

## 6. План доведения до 80% покрытия

1. Ввести интерфейсы зависимостей service-слоя (repository, kinopoisk, youtube).
2. Добавить mock-реализации зависимостей.
3. Написать unit-тесты для `internal/service/service.go`:
   - `Register`, `Login`, `GetProfile`, `UpdateProfile`;
   - `SearchMovies` (cache/db hit, fallback на API, частичные ошибки по страницам);
   - `ListMovies`, `ListPopular` (границы page/size/limit);
   - `GetMovieReviews` (movie-not-found, youtube-fail, success);
   - `Add/Get/Remove watchlist`, `Upsert/Get/Delete ratings`.
4. Добавить позитивные handler-тесты с мок-сервисом (не только валидация, но и success/5xx paths).
5. Включить порог покрытия в CI (например fail ниже 80%).

## 7. Процедура расширения тестового набора (пример)

Пример: добавлен новый метод `GetRecommendations(userID)` в service.

Шаги:
1. Добавить unit-тесты service-метода на классы эквивалентности:
   - пользователь с историей оценок,
   - пользователь без оценок,
   - ошибка репозитория.
2. Добавить unit-тест handler endpoint `/movies/recommendations`:
   - success 200,
   - unauthorized 401,
   - invalid payload/query 400,
   - internal error 500.
3. Зафиксировать новую ветку в отчёте покрытия.
4. Обновить integration/e2e сценарии (если endpoint участвует в пользовательском потоке).
5. Убедиться, что CI job проходит автоматически на push/PR.

## 8. Раздел для вклада участников команды

### Участник 1 (заполнено)
- ФИО:
- Вклад в unit-тесты:
- Файлы:
- Количество тестов:

### Участник 2 (заполнить)
- ФИО:
- Вклад в unit-тесты:
- Файлы:
- Количество тестов:

### Участник 3 (заполнить)
- ФИО:
- Вклад в unit-тесты:
- Файлы:
- Количество тестов:

