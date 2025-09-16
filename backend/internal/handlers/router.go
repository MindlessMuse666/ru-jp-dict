package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Конфигурация Swagger UI
type SwaggerConfig struct {
	FilePath string // Абсолютный путь к файлу swagger.yaml
	URLPath  string // URL-путь для доступа к спецификации
}

// Создает и настраивает маршрутизатор
func SetupRouter(repo *repository.VocabularyRepo, basePath string) *chi.Mux {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Конфигурация Swagger
	swaggerConfig := SwaggerConfig{
		FilePath: filepath.Join(basePath, "docs", "swagger.yaml"),
		URLPath:  "/swagger/openapi.yaml",
	}

	// Настройка маршрутов
	setupSwaggerRoutes(router, swaggerConfig)
	setupAPIRoutes(router, repo)

	return router
}

// Настраивает марщруты Swagger UI
func setupSwaggerRoutes(router *chi.Mux, config SwaggerConfig) {
	// Обслуживаем спецификацию OpenApi
	router.Get(config.URLPath, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, config.FilePath)
	})

	// Настраиваем Swagger UI
	swaggerHandler := httpSwagger.Handler(
		httpSwagger.URL(config.URLPath),
	)

	// Обрабатываем все запросы к Swagger UI
	router.Get("/swagger/*", func(w http.ResponseWriter, r *http.Request) {
		// Редирект запроса к /swagger/doc.json на наш YAML
		if r.URL.Path == "/swagger/doc.json" {
			http.ServeFile(w, r, config.FilePath)
			return
		}
		swaggerHandler.ServeHTTP(w, r)
	})
}

// Настраивает маршруты API
func setupAPIRoutes(router *chi.Mux, repo *repository.VocabularyRepo) {
	vocabHandler := NewVocabularyHandler(repo)

	router.Route("/api/v1", func(r chi.Router) {
		router.Route("/words", func(r chi.Router) {
			router.Get("/", vocabHandler.GetWords)
			router.Post("/", vocabHandler.CreateWord)
			router.Put("/{id}", vocabHandler.UpdateWord)
			router.Delete("/{id}", vocabHandler.DeleteWord)
		})
	})
}
