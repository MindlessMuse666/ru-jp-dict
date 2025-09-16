package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/config"
	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/database"
	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/handlers"
	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/kafka"
	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/repository"
)

func main() {
	// Инициализация БД
	db, err := database.InitDB("./vocabulary.db")
	if err != nil {
		log.Fatal("failed to initialize database: ", err)
	}
	defer db.Close()

	// Инициализация Kafka
	kafkaConfig := config.NewKafkaConfig()
	producer := kafka.NewProducer(kafkaConfig.Broker, kafkaConfig.Topic)
	defer producer.Close()

	// Путь к корню
	basePath := getRootPath()

	// Инициализация репо
	repo := repository.NewVocabularyRepo(db)

	// Настройка HTTP-роутинга
	router := handlers.SetupRouter(repo, producer, basePath)

	// Запуск сервера
	log.Println("server run on: http://localhost:8080")
	log.Println("swagger-ui run on: http://localhost:8080/swagger/index.html")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Получает абсолютный путь к корню проекта
func getRootPath() string {
	// Попытка найти путь относительно текущей директории
	if _, err := os.Stat("./backend/docs"); err == nil {
		return "./backend"
	}
	if _, err := os.Stat("./docs"); err == nil {
		return "."
	}

	// Получение пути к исполняемому файлу
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal("failed to get executable path: ", err)
	}
	return filepath.Dir(execPath)
}
