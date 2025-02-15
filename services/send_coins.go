package services

import (
	"avito-tech-winter-2025/db"
	"avito-tech-winter-2025/dto"
	"database/sql"
	"fmt"
)

type transactionResult struct {
	Err *dto.Error
}

func HandleSendCoinsRequest(token string, request dto.SendCoinsRequest) (*dto.SendCoinsResponse, *dto.Error) {
	userId, err := ExtractJWT(token)
	if err != nil {
		return nil, err
	}

	results := make(chan transactionResult, 3)
	defer close(results)

	var userCoins int
	go func() {
		err := db.DB.QueryRow("SELECT coins FROM employees WHERE id = $1", userId).Scan(&userCoins)
		if err == sql.ErrNoRows {
			results <- transactionResult{Err: &dto.Error{
				Code:       "USER_NOT_FOUND",
				Message:    "Пользователь не найден",
				StatusCode: 404,
			}}
			return
		} else if err != nil {
			results <- transactionResult{Err: &dto.Error{
				Code:       "DB_ERROR",
				Message:    "Ошибка при запросе баланса пользователя",
				StatusCode: 500,
			}}
			return
		}
		results <- transactionResult{}
	}()

	var receiverId int
	go func() {
		err := db.DB.QueryRow("SELECT id FROM employees WHERE username = $1", request.ReceiverUsername).Scan(&receiverId)
		if err == sql.ErrNoRows {
			results <- transactionResult{Err: &dto.Error{
				Code:       "RECEIVER_NOT_FOUND",
				Message:    "Получатель не найден",
				StatusCode: 404,
			}}
			return
		} else if err != nil {
			results <- transactionResult{Err: &dto.Error{
				Code:       "DB_ERROR",
				Message:    "Ошибка при поиске получателя",
				StatusCode: 500,
			}}
			return
		}
		results <- transactionResult{}
	}()

	for i := 0; i < 2; i++ {
		res := <-results
		if res.Err != nil {
			return nil, res.Err
		}
	}

	if userCoins < request.Amount {
		return nil, &dto.Error{
			Code:       "INSUFFICIENT_COINS",
			Message:    "Недостаточно монет",
			StatusCode: 400,
		}
	}

	if receiverId == userId {
		return nil, &dto.Error{
			Code:       "SELF_TRANSFER_NOT_ALLOWED",
			Message:    "Нельзя перевести монеты самому себе",
			StatusCode: 400,
		}
	}

	tx, error := db.DB.Begin()
	if error != nil {
		return nil, &dto.Error{
			Code:       "DB_ERROR",
			Message:    "Ошибка при начале транзакции",
			StatusCode: 500,
		}
	}

	updates := make(chan transactionResult, 3)
	defer close(updates)

	go func() {
		_, err := tx.Exec("UPDATE employees SET coins = coins - $1 WHERE id = $2", request.Amount, userId)
		if err != nil {
			tx.Rollback()
			updates <- transactionResult{Err: &dto.Error{
				Code:       "COINS_DEDUCTION_ERROR",
				Message:    "Ошибка при снятии монет",
				StatusCode: 500,
			}}
			return
		}
		updates <- transactionResult{}
	}()

	go func() {
		_, err := tx.Exec("UPDATE employees SET coins = coins + $1 WHERE id = $2", request.Amount, receiverId)
		if err != nil {
			tx.Rollback()
			updates <- transactionResult{Err: &dto.Error{
				Code:       "COINS_ADDITION_ERROR",
				Message:    "Ошибка при начислении монет",
				StatusCode: 500,
			}}
			return
		}
		updates <- transactionResult{}
	}()

	go func() {
		_, err := tx.Exec("INSERT INTO transactions (sender_id, receiver_id, amount) VALUES ($1, $2, $3)", userId, receiverId, request.Amount)
		if err != nil {
			tx.Rollback()
			updates <- transactionResult{Err: &dto.Error{
				Code:       "TRANSACTION_RECORD_ERROR",
				Message:    "Ошибка при записи транзакции",
				StatusCode: 500,
			}}
			return
		}
		updates <- transactionResult{}
	}()

	for i := 0; i < 3; i++ {
		res := <-updates
		if res.Err != nil {
			return nil, res.Err
		}
	}

	if err := tx.Commit(); err != nil {
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
