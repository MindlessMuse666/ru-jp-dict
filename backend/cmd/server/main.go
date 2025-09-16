package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

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

	// Абсолютный путь к корню
	basePath := getRootPath()

	// Создание репозитория
	repo := repository.NewVocabularyRepo(db)

	// Настройка HTTP-роутинга
	router := handlers.SetupRouter(repo, basePath)

	// Запуск сервера
	log.Println("start server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Получает абсолютный путь к корню проекта
func getRootPath() string {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal("failed to get executable path: ", err)
	}

	return filepath.Dir(execPath)
}
