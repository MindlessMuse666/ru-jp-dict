package handlers

import (
	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

/* Настраивает HTTP-маршруты и возвращает роутер */
func SetupRouter(repo *repository.VocabularyRepo) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

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
