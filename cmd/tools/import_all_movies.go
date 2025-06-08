package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/AlexKeyyyy/movies-picker/internal/models"
	"github.com/AlexKeyyyy/movies-picker/internal/repository"
	"github.com/AlexKeyyyy/movies-picker/pkg/kinopoisk"
	"github.com/joho/godotenv"
)

func main() {
	// Загрузить .env
	_ = godotenv.Load()

	// Подключиться к БД
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL is not set")
	}
	repo, err := repository.NewRepo(dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	// Создать клиента Kinopoisk
	apiKey := os.Getenv("KINOPOISK_API_KEY")
	if apiKey == "" {
		log.Fatal("KINOPOISK_API_KEY is not set")
	}
	client := kinopoisk.NewClient(apiKey)

	// Первая страница — узнаём общее число страниц
	items, totalPages, err := client.GetPopularAll(1)
	if err != nil {
		log.Fatalf("failed to fetch page 1: %v", err)
	}
	log.Printf("Total pages: %d", totalPages)
	upsertFilms(repo, items)

	// Остальные страницы
	for page := 2; page <= totalPages; page++ {
		items, _, err := client.GetPopularAll(page)
		if err != nil {
			log.Printf("failed to fetch page %d: %v", page, err)
			continue
		}
		upsertFilms(repo, items)
	}

	fmt.Println("Import popular movies completed.")
}

func upsertFilms(repo *repository.Repo, films []kinopoisk.Film) {
	for _, f := range films {
		// конвертация года из json.Number
		yearInt, err := strconv.Atoi(f.Year.String())
		if err != nil {
			yearInt = 0
		}
		movie := &models.Movie{
			ID:          f.KinopoiskID,
			Title:       f.NameRu,
			Year:        yearInt,
			Description: f.Description,
			PosterURL:   f.PosterURL,
		}
		if err := repo.UpsertMovie(movie); err != nil {
			log.Printf("upsert failed %d: %v", movie.ID, err)
		} else {
			log.Printf("upserted: %s (%d)", movie.Title, movie.Year)
		}
	}
}
