package handlers

import (
	"avito-tech-winter-2025/services"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func BuyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST-запросы разрешены", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Токен не найден", http.StatusUnauthorized)
		return
	}
	token := authHeader[7:]

	itemName := strings.TrimPrefix(r.URL.Path, "/api/buy/")

	if itemName == "" {
		http.Error(w, "Имя товара не указано", http.StatusBadRequest)
		return
	}

	buyResponse, err := services.HandleBuyRequest(token, itemName)
	if err != nil {
		log.Printf("Ошибка при обработке запроса: %s", err.Message)
		http.Error(w, err.Message, err.StatusCode)
		return
	}

	log.Printf("Успешная покупка товара: %s", itemName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(buyResponse)
}
