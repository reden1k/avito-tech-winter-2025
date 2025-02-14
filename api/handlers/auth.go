package handlers

import (
	"avito-tech-winter-2025/dto"
	"avito-tech-winter-2025/services"
	"encoding/json"
	"net/http"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST-запросы разрешены", http.StatusMethodNotAllowed)
		return
	}

	var req dto.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		authErr := dto.Error{
			Code:    "INVALID_REQUEST",
			Message: "Необходимо передать имя пользователя и пароль",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(authErr)
		return
	}

	authResponse, authErr := services.HandleAuthRequest(req)
	if authErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(authErr.StatusCode)
		json.NewEncoder(w).Encode(authErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResponse)
}
