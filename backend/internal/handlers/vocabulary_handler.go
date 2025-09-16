package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/models"
	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/repository"
	"github.com/go-chi/chi/v5"
)

// Структура репо, которая предоставляет методы-обработчики
type VocabularyHandler struct {
	repo *repository.VocabularyRepo
}

func NewVocabularyHandler(repo *repository.VocabularyRepo) *VocabularyHandler {
	return &VocabularyHandler{repo: repo}
}

// GET /api/v1/words - возвращает все слова
func (h *VocabularyHandler) GetWords(w http.ResponseWriter, r *http.Request) {
	words, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(words)
}

// POST /api/v1/words - создает новое слово
func (h *VocabularyHandler) CreateWord(w http.ResponseWriter, r *http.Request) {
	var word models.Vocabulary

	// мап json в структуру
	if err := json.NewDecoder(r.Body).Decode(&word); err != nil {
		http.Error(w, `{"error": "invalid json"}`, http.StatusBadRequest)
		return
	}

	// валидация: обязательные поля не пустые
	if word.Russian == "" || word.Japanese == "" {
		http.Error(w, `{"error": "fields 'russian' and 'japanese' are required"}`, http.StatusBadRequest)
		return
	}

	// передача данных в репо для создания записи в БД
	id, err := h.repo.Create(word)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// формирование ответа (возвращаем сообщение с ID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	response := map[string]any{
		"id":      id,
		"message": "word created successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// PATCH /appi/v1/words/{id} - частично обновляет слово по ID
func (h *VocabularyHandler) PatchWord(w http.ResponseWriter, r *http.Request) {
	// Извлечение ID из URL-параметра
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, `{"error": "invalid id"}`, http.StatusBadRequest)
		return
	}

	// Декодинг JSON в map вместо структуры
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, `{"error": "invalid json"}`, http.StatusBadRequest)
		return
	}

	// Валидация: переданы только разрешенные поля
	allowedFields := map[string]bool{
		"russian":  true,
		"japanese": true,
		"onyomi":   true,
		"kunoymi":  true,
	}

	for field := range updates {
		if !allowedFields[field] {
			http.Error(w, `{"error": "field '`+field+`' is not allowed}`, http.StatusBadRequest)
			return
		}
	}

	// Валидация: переданные обязательные поля - не пустые
	if russian, ok := updates["russian"]; ok && russian == "" {
		http.Error(w, `{"error": "russian cannot be empty"}`, http.StatusBadRequest)
		return
	}
	if japanese, ok := updates["japanese"]; ok && japanese == "" {
		http.Error(w, `{"error": "japanese cannot be empty"}`, http.StatusBadRequest)
		return
	}

	// Передача данных в репо
	err = h.repo.PartialUpdate(id, updates)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "word partially updated successfully"}
	json.NewEncoder(w).Encode(response)
}

// PUT /appi/v1/words/{id} - полностью обновляет слово по ID
func (h *VocabularyHandler) PutWord(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL-параметра
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, `{"error": "invalid id"}`, http.StatusBadRequest)
		return
	}

	var word models.Vocabulary
	if err := json.NewDecoder(r.Body).Decode(&word); err != nil {
		http.Error(w, `{"error": "invalid json}`, http.StatusBadRequest)
		return
	}

	err = h.repo.Update(id, word)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	responce := map[string]string{"message": "word updated successfully"}
	json.NewEncoder(w).Encode(responce)
}

// DELETE /api/v1/words/{id} - удаляет слово по ID
func (h *VocabularyHandler) DeleteWord(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, `{"error": "invalid id"}`, http.StatusBadRequest)
		return
	}

	err = h.repo.Delete(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	responce := map[string]string{"message": "word deleted successfully"}
	json.NewEncoder(w).Encode(responce)
}
