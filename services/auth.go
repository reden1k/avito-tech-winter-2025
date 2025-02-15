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

	resultCh := make(chan *dto.AuthResponse)
	errorCh := make(chan *dto.Error)

	go func() {
		var id int
		var hashedPassword string
		err := db.DB.QueryRow("SELECT id, password_hash FROM employees WHERE username = $1", req.Username).
			Scan(&id, &hashedPassword)

		if err == sql.ErrNoRows {
			response, err := createUser(req.Username, req.Password)
			if err != nil {
				errorCh <- err
			} else {
				resultCh <- response
			}
			return
		} else if err != nil {
			errorCh <- &dto.Error{
				Code:       "DB_ERROR",
				Message:    "Ошибка работы с базой данных",
				StatusCode: 500,
			}
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)) != nil {
			errorCh <- &dto.Error{
				Code:       "INVALID_PASSWORD",
				Message:    "Неверный пароль",
				StatusCode: 401,
			}
			return
		}

		tokenString, err := generateJWT(id)
		if err != nil {
			errorCh <- &dto.Error{
				Code:       "TOKEN_GENERATION_ERROR",
				Message:    "Ошибка генерации токена",
				StatusCode: 500,
			}
			return
		}

		resultCh <- &dto.AuthResponse{Token: tokenString}
	}()

	select {
	case res := <-resultCh:
		return res, nil
	case err := <-errorCh:
		return nil, err
	}
}

func createUser(username, password string) (*dto.AuthResponse, *dto.Error) {
	hashCh := make(chan string)
	errorCh := make(chan *dto.Error)

	go func() {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			errorCh <- &dto.Error{
				Code:       "PASSWORD_HASH_ERROR",
				Message:    "Ошибка хеширования пароля",
				StatusCode: 500,
			}
			return
		}
		hashCh <- string(hashedPassword)
	}()

	select {
	case hashedPassword := <-hashCh:
		var id int
		err := db.DB.QueryRow(
			"INSERT INTO employees (username, password_hash, coins) VALUES ($1, $2, $3) RETURNING id",
			username, hashedPassword, 1000,
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
	case err := <-errorCh:
		return nil, err
	}
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
