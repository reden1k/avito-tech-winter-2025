package services

import (
	"avito-tech-winter-2025/dto"
	"avito-tech-winter-2025/models"
	"database/sql"
	"fmt"
)

func HandleSendCoinsRequest(token string, request dto.SendCoinsRequest) (*dto.SendCoinsResponse, *dto.Error) {
	userId, error := ExtractJWT(token)
	if error != nil {
		return nil, error
	}

	var userCoins int
	err := models.DB.QueryRow("SELECT coins FROM employees WHERE id = $1", int(userId)).Scan(&userCoins)
	if err == sql.ErrNoRows {
		return nil, &dto.Error{
			Code:       "USER_NOT_FOUND",
			Message:    "Пользователь не найден",
			StatusCode: 404,
		}
	} else if err != nil {
		return nil, &dto.Error{
			Code:       "DB_ERROR",
			Message:    "Ошибка при запросе баланса пользователя",
			StatusCode: 500,
		}
	}

	if userCoins < request.Amount {
		return nil, &dto.Error{
			Code:       "INSUFFICIENT_COINS",
			Message:    "Недостаточно монет",
			StatusCode: 400,
		}
	}

	var receiverId int
	err = models.DB.QueryRow("SELECT id FROM employees WHERE username = $1", request.ReceiverUsername).Scan(&receiverId)
	if err == sql.ErrNoRows {
		return nil, &dto.Error{
			Code:       "RECEIVER_NOT_FOUND",
			Message:    "Получатель не найден",
			StatusCode: 404,
		}
	} else if err != nil {
		return nil, &dto.Error{
			Code:       "DB_ERROR",
			Message:    "Ошибка при поиске получателя",
			StatusCode: 500,
		}
	}

	if receiverId == int(userId) {
		return nil, &dto.Error{
			Code:       "SELF_TRANSFER_NOT_ALLOWED",
			Message:    "Нельзя перевести монеты самому себе",
			StatusCode: 400,
		}
	}

	tx, err := models.DB.Begin()
	if err != nil {
		return nil, &dto.Error{
			Code:       "DB_ERROR",
			Message:    "Ошибка при начале транзакции",
			StatusCode: 500,
		}
	}

	_, err = tx.Exec("UPDATE employees SET coins = coins - $1 WHERE id = $2", request.Amount, int(userId))
	if err != nil {
		tx.Rollback()
		return nil, &dto.Error{
			Code:       "COINS_DEDUCTION_ERROR",
			Message:    "Ошибка при снятии монет",
			StatusCode: 500,
		}
	}

	_, err = tx.Exec("UPDATE employees SET coins = coins + $1 WHERE id = $2", request.Amount, receiverId)
	if err != nil {
		tx.Rollback()
		return nil, &dto.Error{
			Code:       "COINS_ADDITION_ERROR",
			Message:    "Ошибка при начислении монет",
			StatusCode: 500,
		}
	}

	_, err = tx.Exec("INSERT INTO transactions (sender_id, receiver_id, amount) VALUES ($1, $2, $3)", int(userId), receiverId, request.Amount)
	if err != nil {
		tx.Rollback()
		return nil, &dto.Error{
			Code:       "TRANSACTION_RECORD_ERROR",
			Message:    "Ошибка при записи транзакции",
			StatusCode: 500,
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, &dto.Error{
			Code:       "TRANSACTION_COMMIT_ERROR",
			Message:    "Ошибка при завершении транзакции",
			StatusCode: 500,
		}
	}

	return &dto.SendCoinsResponse{
		Message: fmt.Sprintf("Успешно отправлено %d монет пользователю %s", request.Amount, request.ReceiverUsername),
	}, nil
}
