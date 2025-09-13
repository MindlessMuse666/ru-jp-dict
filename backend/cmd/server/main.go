package main

import (
	"fmt"
	"log"
	"net/http"
)

// Обработчик для главной страницы
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Привет! Это бэкенд ru-jp-dict!</h1><p>Сервер работает.</p>")
}

func main() {
	// Конфигурируем маршруты
	http.HandleFunc("/", homeHandler)

	// Запускаем сервер на порту 8080
	fmt.Println("Запуск сервера на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
