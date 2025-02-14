package handlers

import (
	"avito-tech-winter-2025/services"
	"encoding/json"
	"net/http"
	"strings"
)

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Только GET-запросы разрешены", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Токен не найден", http.StatusUnauthorized)
		return
	}
	token := authHeader[7:]

	infoResponse, err := services.HandleInfoRequest(token)
	if err != nil {
		http.Error(w, err.Message, err.StatusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(infoResponse)
}
