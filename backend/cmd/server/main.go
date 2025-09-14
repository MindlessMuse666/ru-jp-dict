package main

import (
	"log"
	"net/http"

	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/database"
	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/handlers"
	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/repository"
)

func main() {
	// Инициализация БД
	db, err := database.InitDB("./vocabulary.db")
	if err != nil {
		log.Fatal("failed to initialize database: ", err)
	}
	defer db.Close()

	// Создание репозитория
	repo := repository.NewVocabularyRepo(db)

	/* // МОК-ДАННЫЕ ДЛЯ ТЕСТА
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
	} */

	// Настройка HTTP-роутинга
	router := handlers.SetupRouter(repo)

	// Запуск сервера
	log.Println("запуск сервера на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
