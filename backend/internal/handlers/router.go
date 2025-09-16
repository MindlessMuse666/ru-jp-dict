package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/kafka"
	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Конфигурация Swagger UI
type SwaggerConfig struct {
	FilePath string // Абсолютный путь к файлу swagger.yaml
}

// Создает и настраивает маршрутизатор
func SetupRouter(repo *repository.VocabularyRepo, producer *kafka.Producer, basePath string) *chi.Mux {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Конфигурация Swagger
	swaggerConfig := SwaggerConfig{
		FilePath: filepath.Join(basePath, "docs", "swagger.yaml"),
	}

	// Обработчик для корневого пути
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})

	// Настройка маршрутов
	setupAPIRoutes(router, repo, producer)
	setupSwaggerRoutes(router, swaggerConfig)

	return router
}

// Настраивает марщруты Swagger UI
func setupSwaggerRoutes(router *chi.Mux, config SwaggerConfig) {
	// Обслуживаем спецификацию OpenApi
	router.Get("/swagger/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		http.ServeFile(w, r, config.FilePath)
	})

	// Настраиваем Swagger UI
	swaggerHandler := httpSwagger.Handler(
		httpSwagger.URL("/swagger/openapi.yaml"),
	)

	// Обрабатываем все запросы к Swagger UI
	router.Get("/swagger/*", swaggerHandler)
}

// Настраивает маршруты API
func setupAPIRoutes(router *chi.Mux, repo *repository.VocabularyRepo, producer *kafka.Producer) {
	vocabHandler := NewVocabularyHandler(repo, producer)

	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/words", func(r chi.Router) {
			r.Get("/", vocabHandler.GetWords)
			r.Post("/", vocabHandler.CreateWord)
			r.Put("/{id}", vocabHandler.PutWord)
			r.Patch("/{id}", vocabHandler.PatchWord)
			r.Delete("/{id}", vocabHandler.DeleteWord)
		})
	})
}
