package handlers

import (
	"avito-tech-winter-2025/dto"
	"avito-tech-winter-2025/services"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func SendCoinsHandler(w http.ResponseWriter, r *http.Request) {
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

	var sendCoinsRequest dto.SendCoinsRequest
	if err := json.NewDecoder(r.Body).Decode(&sendCoinsRequest); err != nil {
		http.Error(w, "Ошибка декодирования запроса", http.StatusBadRequest)
		return
	}

	sendCoinsResponse, err := services.HandleSendCoinsRequest(token, sendCoinsRequest)
	if err != nil {
		log.Printf("Ошибка при переводе монет: %s", err.Message)
		http.Error(w, err.Message, err.StatusCode)
		return
	}

	log.Printf("Успешный перевод %d монет пользователю %s", sendCoinsRequest.Amount, sendCoinsRequest.ReceiverUsername)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sendCoinsResponse)
}
