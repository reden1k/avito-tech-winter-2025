package services

import (
	"avito-tech-winter-2025/db"
	"avito-tech-winter-2025/dto"
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HandleAuthRequest(req dto.AuthRequest) (*dto.AuthResponse, *dto.Error) {
	if req.Username == "" || req.Password == "" {
		return nil, &dto.Error{
			Code:       "MISSING_CREDENTIALS",
			Message:    "Необходимо ввести имя пользователя и пароль",
			StatusCode: 400,
		}
	}

	var id int
	var hashedPassword string
	err := db.DB.QueryRow("SELECT id, password_hash FROM employees WHERE username = $1", req.Username).
		Scan(&id, &hashedPassword)

	if err == sql.ErrNoRows {
		return createUser(req.Username, req.Password)
	} else if err != nil {
		return nil, &dto.Error{
			Code:       "DB_ERROR",
			Message:    "Ошибка работы с базой данных",
			StatusCode: 500,
		}
	}

	// Проверяем пароль
	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)) != nil {
		return nil, &dto.Error{
			Code:       "INVALID_PASSWORD",
			Message:    "Неверный пароль",
			StatusCode: 401,
		}
	}

	// Генерируем JWT
	tokenString, err := generateJWT(id)
	if err != nil {
		return nil, &dto.Error{
			Code:       "TOKEN_GENERATION_ERROR",
			Message:    "Ошибка генерации токена",
			StatusCode: 500,
		}
	}

	return &dto.AuthResponse{Token: tokenString}, nil
}

func createUser(username, password string) (*dto.AuthResponse, *dto.Error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return nil, &dto.Error{
			Code:       "PASSWORD_HASH_ERROR",
			Message:    "Ошибка хеширования пароля",
			StatusCode: 500,
		}
	}

	var id int
	err = db.DB.QueryRow(
		"INSERT INTO employees (username, password_hash, coins) VALUES ($1, $2, $3) RETURNING id",
		username, string(hashedPassword), 1000,
	).Scan(&id)

	if err != nil {
		return nil, &dto.Error{
			Code:       "USER_CREATION_ERROR",
			Message:    "Ошибка создания пользователя",
			StatusCode: 500,
		}
	}

	tokenString, err := generateJWT(id)
	if err != nil {
		return nil, &dto.Error{
			Code:       "TOKEN_GENERATION_ERROR",
			Message:    "Ошибка генерации токена",
			StatusCode: 500,
		}
	}

	return &dto.AuthResponse{Token: tokenString}, nil
}

func generateJWT(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 6).Unix(),
	})

	return token.SignedString([]byte("jwt_secret_token"))
}

func ExtractJWT(token string) (int, *dto.Error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("jwt_secret_token"), nil
	})

	if err != nil {
		return 0, &dto.Error{
			Code:       "INVALID_TOKEN",
			Message:    "Неверный токен",
			StatusCode: 401,
		}
	}

	userId, ok := claims["user_id"].(float64)
	if !ok {
		return 0, &dto.Error{
			Code:       "USER_ID_NOT_FOUND",
			Message:    "Не удалось найти user_id в токене",
			StatusCode: 500,
		}
	}

	return int(userId), nil
}
