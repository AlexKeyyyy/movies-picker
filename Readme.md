# Movie Picker (Movie Tracker)

**Веб-сервис подбора и отслеживания фильмов**  
Позволяет:

- Искать фильмы по названию или ключевым словам (через Kinopoisk API).
- Просматривать детали фильма (описание, год, постер).
- Добавлять фильмы в личный список "Смотреть позже".
- Сохранять собственные рейтинги фильмов.
- Получать ссылки на авторитетные видео-разборы (YouTube Data API).

---

## 📦 Технологический стек

- **Backend:** Go
- **Database:** PostgreSQL
- **Внешние API:**
  - Kinopoisk Unofficial API (фильмы)
  - YouTube Data API v3 (обзоры)
- **Контейнеризация (опционально):** Docker & Docker Compose
- **Тестирование:**
  - Unit-тесты (Go)
  - Интеграционные тесты (httptest)

---

## 🚀 Быстрый старт

1. **Клонировать репозиторий**

   ```bash
   git clone https://github.com/AlexKeyyyy/movies-picker.git
   cd movies-picker
   ```

2. **Создать файл окружения `.env`**

   ```ini
   PORT=8080
   DB_URL=postgres://user:password@localhost:5432/moviedb?sslmode=disable
   JWT_SECRET=your_jwt_secret
   KINOPOISK_API_KEY=your_kinopoisk_api_key
   ```

3. **Установить зависимости**

   ```bash
   go mod tidy
   ```

4. **Запустить локальную БД**

   - Через Docker Compose:

     ```bash
     docker compose up -d postgres
     ```

   - Или вручную (pgAdmin / локальный сервер).

5. **Запустить сервер**

   ```bash
   go run cmd/server/main.go
   ```

---

## 🗄️ База данных

### Создание таблиц (PostgreSQL)

```sql
CREATE TABLE users (
  user_id       SERIAL PRIMARY KEY,
  email         VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE movies (
  movie_id    BIGINT PRIMARY KEY,
  title       VARCHAR(255) NOT NULL,
  year        INT,
  poster_url  TEXT,
  description TEXT,
  last_sync   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE watchlist (
  user_id    INT NOT NULL,
  movie_id   BIGINT NOT NULL,
  added_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(user_id, movie_id),
  FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE,
  FOREIGN KEY(movie_id) REFERENCES movies(movie_id) ON DELETE CASCADE
);

CREATE TABLE ratings (
  user_id  INT NOT NULL,
  movie_id BIGINT NOT NULL,
  rating   SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 10),
  rated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(user_id, movie_id),
  FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE,
  FOREIGN KEY(movie_id) REFERENCES movies(movie_id) ON DELETE CASCADE
);
```

---

## Сборка

На Windows:

1. Создать .env файл, например:
   PORT=8080
   DB_URL=postgres://postgres:alexkoba@postgres:5432/moviedb?sslmode=disable
   JWT_SECRET=hello
   KINOPOISK_API_KEY=bda3897d-a997-48e4-97a8-0d0bd514e7b3
   YOUTUBE_API_KEY=AIzaSyBGnHMN-tDg43gjwwpXcWj1fsjXTH28oQw

2. Запустить ./run.ps1

На Linux:

1. Аналогично
2. Запустить ./run.sh

## Запуск локально

1. Создать .env файл, например:
   PORT=8080
   DB_URL=postgres://postgres:alexkoba@localhost:5432/moviedb?sslmode=disable
   JWT_SECRET=hello
   KINOPOISK_API_KEY=bda3897d-a997-48e4-97a8-0d0bd514e7b3
   YOUTUBE_API_KEY=AIzaSyBGnHMN-tDg43gjwwpXcWj1fsjXTH28oQw

2. go run cmd/server/main.go

Спецификация Swagger находится по адресу, например:
localhost:PORT/docs
