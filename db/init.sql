CREATE TABLE IF NOT EXISTS users (
  user_id       SERIAL PRIMARY KEY,
  email         VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS movies (
  movie_id    BIGINT PRIMARY KEY,
  title       VARCHAR(255) NOT NULL,
  year        INT,
  poster_url  TEXT,
  description TEXT,
  last_sync   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS watchlist (
  user_id   INT NOT NULL,
  movie_id  BIGINT NOT NULL,
  added_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(user_id, movie_id),
  FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE,
  FOREIGN KEY(movie_id) REFERENCES movies(movie_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ratings (
  user_id  INT NOT NULL,
  movie_id BIGINT NOT NULL,
  rating   SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 10),
  rated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(user_id, movie_id),
  FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE,
  FOREIGN KEY(movie_id) REFERENCES movies(movie_id) ON DELETE CASCADE
);
