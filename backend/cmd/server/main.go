package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/database"
	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/models"
	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/repository"
)

func main() {
	// Инициализация БД
	db, err := database.InitDB("./vocabulary.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание репозитория
	repo := repository.NewVocabularyRepo(db)

	// МОК-ДАННЫЕ ДЛЯ ТЕСТА
	id, err := repo.Create(models.Vocabulary{
		Russian:  "книга",
		Japanese: "木",
		Onyomi:   "き",
		Kunyomi:  "モク",
	})

	if err != nil {
		log.Fatal("ошибка при создании слова: ", err)
	}

	fmt.Printf("добавлено слово с ID: %d\n", id)

	words, err := repo.GetAll()
	if err != nil {
		log.Fatal("ошибка при получении слов: ", err)
	}

	fmt.Println("слова в базе:")
	for _, word := range words {
		fmt.Printf(
			"%d: %s - %s (On: %s, Kun: %s)\n",
			word.ID,
			word.Russian,
			word.Japanese,
			word.Onyomi,
			word.Kunyomi,
		)
	}

	// Конфигурируем маршруты
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>Привет! Это бэкенд ru-jp-dict!</h1><p>Сервер работает. База данных подключена.</p>"))
	})

	// Запускаем сервер на порту 8080
	fmt.Println("Запуск сервера на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
