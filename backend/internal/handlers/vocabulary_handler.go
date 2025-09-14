package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/repository"
)

/* Хранит репо и предоставляет методы-обработчики */
type VocabularyHandler struct {
	repo *repository.VocabularyRepo
}

func NewVocabularyHandler(repo *repository.VocabularyRepo) *VocabularyHandler {
	return &VocabularyHandler{repo: repo}
}

/* GET /api/v1/words - возвращает все слова */
func (h *VocabularyHandler) GetWords(w http.ResponseWriter, r *http.Request) {
	words, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(words)
}

func (h *VocabularyHandler) CreateWord(w http.ResponseWriter, r *http.Request) {
	// TODO
}
