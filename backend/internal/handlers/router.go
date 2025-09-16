package handlers

import (
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

/* Настраивает HTTP-маршруты и возвращает роутер */
func SetupRouter(repo *repository.VocabularyRepo) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Абсолютный путь к swagger.yaml
	_, filename, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	swaggerPath := filepath.Join(basepath, "docs", "swagger.yaml")

	// Обслуживание спецификации
	r.Get("/swagger/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, swaggerPath)
	})

	swaggerHandler := httpSwagger.Handler(
		httpSwagger.URL("/swagger/openapi.yaml"),
	)

	r.Get("/swagger/*", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/swagger/doc.json" {
			http.ServeFile(w, r, swaggerPath)
			return
		}
		swaggerHandler.ServeHTTP(w, r)
	})

	vocabHandler := NewVocabularyHandler(repo)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/words", func(r chi.Router) {
			r.Get("/", vocabHandler.GetWords)
			r.Post("/", vocabHandler.CreateWord)
			r.Put("/{id}", vocabHandler.UpdateWord)
			r.Delete("/{id}", vocabHandler.DeleteWord)
		})
	})

	return r
}
