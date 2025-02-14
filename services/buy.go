package services

import (
	"avito-tech-winter-2025/dto"
	"avito-tech-winter-2025/models"
	"database/sql"
	"fmt"
)

func HandleBuyRequest(token string, itemName string) (*dto.BuyResponse, *dto.Error) {
	userId, error := ExtractJWT(token)
	if error != nil {
		return nil, error
	}

	tx, err := models.DB.Begin()
	if err != nil {
		return nil, &dto.Error{
			Code:       "TRANSACTION_ERROR",
			Message:    "Ошибка при создании транзакции",
			StatusCode: 500,
		}
	}
	defer tx.Rollback()

	var itemPrice, userCoins int
	err = tx.QueryRow(`
		SELECT price, coins
		FROM items
		JOIN employees ON employees.id = $1
		WHERE items.name = $2
	`, int(userId), itemName).Scan(&itemPrice, &userCoins)

	if err == sql.ErrNoRows {
		return nil, &dto.Error{
			Code:       "ITEM_NOT_FOUND",
			Message:    "Товар не найден",
			StatusCode: 404,
		}
	} else if err != nil {
		return nil, &dto.Error{
			Code:       "ITEM_RETRIEVAL_ERROR",
			Message:    "Ошибка при извлечении товара или пользователя",
			StatusCode: 500,
		}
	}

	if userCoins < itemPrice {
		return nil, &dto.Error{
			Code:       "INSUFFICIENT_FUNDS",
			Message:    "Недостаточно монет для покупки",
			StatusCode: 400,
		}
	}

	_, err = tx.Exec(`
		INSERT INTO purchases (employee_id, item_id)
		SELECT $1, id FROM items WHERE name = $2
	`, int(userId), itemName)
	if err != nil {
		return nil, &dto.Error{
			Code:       "PURCHASE_ERROR",
			Message:    "Ошибка при регистрации покупки",
			StatusCode: 500,
		}
	}

	_, err = tx.Exec(`
		UPDATE employees SET coins = coins - $1 WHERE id = $2
	`, itemPrice, int(userId))
	if err != nil {
		return nil, &dto.Error{
			Code:       "COINS_DEDUCTION_ERROR",
			Message:    "Ошибка при снятии монет с пользователя",
			StatusCode: 500,
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
