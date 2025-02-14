package main

import (
	"avito-tech-winter-2025/api/handlers"
	"avito-tech-winter-2025/models"
	"log"
	"net/http"
)

func main() {
	models.InitDB()

	http.HandleFunc("/api/auth", handlers.AuthHandler)
	http.HandleFunc("/api/info", handlers.InfoHandler)
	http.HandleFunc("/api/buy/", handlers.BuyHandler)
	http.HandleFunc("/api/sendCoin", handlers.SendCoinsHandler)

	log.Println("Сервер запущен на порту 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
