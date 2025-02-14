package services

import (
	"avito-tech-winter-2025/dto"
	"avito-tech-winter-2025/models"
	"database/sql"
)

func HandleInfoRequest(token string) (*dto.InfoResponse, *dto.Error) {
	userId, error := ExtractJWT(token)
	if error != nil {
		return nil, error
	}

	var user models.Employee
	err := models.DB.QueryRow("SELECT id, coins FROM employees WHERE id = $1", int(userId)).Scan(&user.ID, &user.Coins)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &dto.Error{
				Code:       "USER_NOT_FOUND",
				Message:    "Пользователь не найден",
				StatusCode: 404,
			}
		}
		return nil, &dto.Error{
			Code:       "DB_ERROR",
			Message:    "Ошибка базы данных",
			StatusCode: 500,
		}
	}

	var inventory []dto.InventoryItem
	rows, err := models.DB.Query(`
        SELECT items.name, COUNT(*) 
        FROM purchases 
        JOIN items ON purchases.item_id = items.id 
        WHERE purchases.employee_id = $1 
        GROUP BY items.name`, user.ID)
	if err != nil {
		return nil, &dto.Error{
			Code:       "ITEMS_RETRIEVAL_ERROR",
			Message:    "Ошибка при извлечении инвентаря",
			StatusCode: 500,
		}
	}
	defer rows.Close()

	for rows.Next() {
		var item dto.InventoryItem
		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			return nil, &dto.Error{
				Code:       "ITEMS_SCAN_ERROR",
				Message:    "Ошибка при обработке инвентаря",
				StatusCode: 500,
			}
		}
		inventory = append(inventory, item)
	}

	var received []dto.CoinTransaction
	var sent []dto.CoinTransaction
	rows, err = models.DB.Query(`
        SELECT sender_id, receiver_id, amount 
        FROM transactions 
        WHERE sender_id = $1 OR receiver_id = $1`, user.ID)
	if err != nil {
		return nil, &dto.Error{
			Code:       "TRANSACTIONS_RETRIEVAL_ERROR",
			Message:    "Ошибка при извлечении транзакций",
			StatusCode: 500,
		}
	}
	defer rows.Close()

	for rows.Next() {
		var t dto.CoinTransaction
		var senderId, receiverId int
		if err := rows.Scan(&senderId, &receiverId, &t.Amount); err != nil {
			return nil, &dto.Error{
				Code:       "TRANSACTIONS_SCAN_ERROR",
				Message:    "Ошибка при обработке транзакций",
				StatusCode: 500,
			}
		}

		var receiverUsername string
		err := models.DB.QueryRow("SELECT username FROM employees WHERE id = $1", receiverId).Scan(&receiverUsername)
		if err != nil {
			return nil, &dto.Error{
				Code:       "RECEIVER_USERNAME_ERROR",
				Message:    "Ошибка при получении имени пользователя получателя",
				StatusCode: 500,
			}
		}

		if senderId == user.ID {
			t.ToUser = receiverUsername
			sent = append(sent, t)
		} else {
			var senderUsername string
			err := models.DB.QueryRow("SELECT username FROM employees WHERE id = $1", senderId).Scan(&senderUsername)
			if err != nil {
				return nil, &dto.Error{
					Code:       "SENDER_USERNAME_ERROR",
					Message:    "Ошибка при получении имени пользователя отправителя",
					StatusCode: 500,
				}
			}
			t.FromUser = senderUsername
			received = append(received, t)
		}
	}

	infoResponse := &dto.InfoResponse{
		Coins:     user.Coins,
		Inventory: inventory,
		CoinHistory: dto.CoinHistory{
			Received: received,
			Sent:     sent,
		},
	}

	return infoResponse, nil
}
