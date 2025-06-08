// package main

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/AlexKeyyyy/movies-picker/internal/models"
// 	"github.com/AlexKeyyyy/movies-picker/internal/repository"
// 	"github.com/AlexKeyyyy/movies-picker/pkg/kinopoisk"
// 	"github.com/joho/godotenv"
// )

// func main() {
// 	if err := godotenv.Load(); err != nil {
// 		log.Println("No .env file found")
// 	}

// 	dsn := os.Getenv("DB_URL")
// 	if dsn == "" {
// 		log.Fatal("DB_URL is not set")
// 	}

// 	repo, err := repository.NewRepo(dsn)
// 	if err != nil {
// 		log.Fatalf("failed to connect to database: %v", err)
// 	}

// 	apiKey := os.Getenv("KINOPOISK_API_KEY")
// 	if apiKey == "" {
// 		log.Fatal("KINOPOISK_API_KEY is not set")
// 	}

// 	client := kinopoisk.NewClient(apiKey)

// 	for page := 1; page <= 13; page++ {
// 		films, err := client.GetTop250Films(page)
// 		if err != nil {
// 			log.Printf("failed to get films on page %d: %v", page, err)
// 			continue
// 		}

// 		for _, f := range films {
// 			// Преобразуем год из string в int
// 			yearInt, err := f.Year.Int64()
// 			if err != nil {
// 				// если не получилось, можно пропустить или установить 0
// 				yearInt = 0
// 			}

// 			movie := &models.Movie{
// 				ID:          f.KinopoiskID,
// 				Title:       f.NameRu,
// 				Year:        int(yearInt),
// 				Description: f.Description,
// 				PosterURL:   f.PosterURL,
// 			}
// 			log.Printf(movie.Title)
// 			if err := repo.UpsertMovie(movie); err != nil {
// 				log.Printf("failed to upsert movie %d: %v", f.KinopoiskID, err)
// 			} else {
// 				log.Printf("upserted movie: %s (%d)", movie.Title, movie.Year)
// 			}
// 		}
// 	}

// 	fmt.Println("Import completed.")
// }
