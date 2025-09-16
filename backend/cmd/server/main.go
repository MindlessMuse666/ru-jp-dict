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

	// Настройка HTTP-роутинга
	router := handlers.SetupRouter(repo)

	// Запуск сервера
	log.Println("запуск сервера на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
