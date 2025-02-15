package services

import (
	"avito-tech-winter-2025/db"
	"avito-tech-winter-2025/dto"
	"database/sql"
	"fmt"
)

type buyResult struct {
	itemPrice int
	userCoins int
	err       *dto.Error
}

func HandleBuyRequest(token string, itemName string) (*dto.BuyResponse, *dto.Error) {
	userId, err := ExtractJWT(token)
	if err != nil {
		return nil, err
	}

	resultCh := make(chan buyResult)

	go func() {
		var itemPrice, userCoins int
		err := db.DB.QueryRow(`
			SELECT items.price, employees.coins 
			FROM items 
			JOIN employees ON employees.id = $1 
			WHERE items.name = $2`, int(userId), itemName).Scan(&itemPrice, &userCoins)

		if err == sql.ErrNoRows {
			resultCh <- buyResult{err: &dto.Error{
				Code:       "ITEM_NOT_FOUND",
				Message:    "Товар не найден",
				StatusCode: 404,
			}}
		} else if err != nil {
			resultCh <- buyResult{err: &dto.Error{
				Code:       "ITEM_RETRIEVAL_ERROR",
				Message:    "Ошибка при извлечении товара или пользователя",
				StatusCode: 500,
			}}
		} else {
			resultCh <- buyResult{itemPrice: itemPrice, userCoins: userCoins}
		}
	}()

	result := <-resultCh
	if result.err != nil {
		return nil, result.err
	}

	if result.userCoins < result.itemPrice {
		return nil, &dto.Error{
			Code:       "INSUFFICIENT_FUNDS",
			Message:    "Недостаточно монет для покупки",
			StatusCode: 400,
		}
	}

	tx, error := db.DB.Begin()
	if error != nil {
		return nil, &dto.Error{
			Code:       "TRANSACTION_ERROR",
			Message:    "Ошибка при создании транзакции",
			StatusCode: 500,
		}
	}
	defer tx.Rollback()

	errCh := make(chan *dto.Error, 2)

	done := make(chan bool)

	go func() {
		_, err := tx.Exec(`INSERT INTO purchases (employee_id, item_id) SELECT $1, id FROM items WHERE name = $2`, int(userId), itemName)
		if err != nil {
			errCh <- &dto.Error{
				Code:       "PURCHASE_ERROR",
				Message:    "Ошибка при регистрации покупки",
				StatusCode: 500,
			}
			return
		}
		done <- true
	}()

	go func() {
		_, err := tx.Exec(`UPDATE employees SET coins = coins - $1 WHERE id = $2`, result.itemPrice, int(userId))
		if err != nil {
			errCh <- &dto.Error{
				Code:       "COINS_DEDUCTION_ERROR",
				Message:    "Ошибка при снятии монет с пользователя",
				StatusCode: 500,
			}
			return
		}
		done <- true
	}()

	for i := 0; i < 2; i++ {
		select {
		case <-done:
		case err := <-errCh:
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, &dto.Error{
			Code:       "TRANSACTION_COMMIT_ERROR",
			Message:    "Ошибка при завершении транзакции",
			StatusCode: 500,
		}
	}

	return &dto.BuyResponse{
		Message: fmt.Sprintf("Товар '%s' успешно куплен", itemName),
	}, nil
}
