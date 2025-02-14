package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"avito-tech-winter-2025/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("jwt_secret_token")

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST-запросы разрешены", http.StatusMethodNotAllowed)
		return
	}

	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	var id int
	var hashedPassword string
	err := models.DB.QueryRow("SELECT id, password_hash FROM employees WHERE username = $1", req.Username).
		Scan(&id, &hashedPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if hashErr != nil {
				http.Error(w, "Ошибка хеширования пароля", http.StatusInternalServerError)
				return
			}

			err = models.DB.QueryRow(
				"INSERT INTO employees (username, password_hash, coins) VALUES ($1, $2, $3) RETURNING id",
				req.Username, string(hashedPassword), 1000,
			).Scan(&id)

			if err != nil {
				http.Error(w, "Ошибка создания пользователя", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Ошибка при проверке пользователя", http.StatusInternalServerError)
			return
		}
	} else {
		if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)) != nil {
			http.Error(w, "Неверный пароль", http.StatusUnauthorized)
			return
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(time.Hour * 6).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	response := models.AuthResponse{Token: tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
