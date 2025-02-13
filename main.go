package main

import (
	"avito-tech-winter-2025/api/handlers"
	"avito-tech-winter-2025/models"
	"log"
	"net/http"
)

func main() {
	models.InitDB()

	http.HandleFunc("/auth", handlers.AuthHandler)

	log.Println("Сервер запущен на порту 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
