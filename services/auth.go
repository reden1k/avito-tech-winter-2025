package services

import (
	"avito-tech-winter-2025/dto"
	"avito-tech-winter-2025/models"
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

	tx, err := models.DB.Begin()
	if err != nil {
		return nil, &dto.Error{
			Code:       "DB_TRANSACTION_ERROR",
			Message:    "Ошибка при создании транзакции",
			StatusCode: 500,
		}
	}
	defer tx.Rollback()

	var id int
	var hashedPassword string
	err = tx.QueryRow("SELECT id, password_hash FROM employees WHERE username = $1", req.Username).
		Scan(&id, &hashedPassword)

	if err == sql.ErrNoRows {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, &dto.Error{
				Code:       "PASSWORD_HASH_ERROR",
				Message:    "Ошибка хеширования пароля",
				StatusCode: 500,
			}
		}

		err = tx.QueryRow(
			"INSERT INTO employees (username, password_hash, coins) VALUES ($1, $2, $3) RETURNING id",
			req.Username, string(hashedPassword), 1000,
		).Scan(&id)

		if err != nil {
			return nil, &dto.Error{
				Code:       "USER_CREATION_ERROR",
				Message:    "Ошибка создания пользователя",
				StatusCode: 500,
			}
		}

		// Генерация токена
		tokenString, err := generateJWT(id)
		if err != nil {
			return nil, &dto.Error{
				Code:       "TOKEN_GENERATION_ERROR",
				Message:    "Ошибка генерации токена",
				StatusCode: 500,
			}
		}

		if err := tx.Commit(); err != nil {
			return nil, &dto.Error{
				Code:       "DB_COMMIT_ERROR",
				Message:    "Ошибка при подтверждении транзакции",
				StatusCode: 500,
			}
		}

		return &dto.AuthResponse{Token: tokenString}, nil
	} else if err != nil {
		return nil, &dto.Error{
			Code:       "DB_ERROR",
			Message:    "Ошибка работы с базой данных",
			StatusCode: 500,
		}
	}

	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)) != nil {
		return nil, &dto.Error{
			Code:       "INVALID_PASSWORD",
			Message:    "Неверный пароль",
			StatusCode: 401,
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

	if err := tx.Commit(); err != nil {
		return nil, &dto.Error{
			Code:       "DB_COMMIT_ERROR",
			Message:    "Ошибка при подтверждении транзакции",
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
