package services

import (
	"avito-tech-winter-2025/db"
	"avito-tech-winter-2025/dto"
	"avito-tech-winter-2025/models"
	"database/sql"
)

func HandleInfoRequest(token string) (*dto.InfoResponse, *dto.Error) {
	userId, err := ExtractJWT(token)
	if err != nil {
		return nil, err
	}

	userCh := make(chan *models.Employee)
	invCh := make(chan []dto.InventoryItem)
	receivedCh := make(chan []dto.CoinTransaction)
	sentCh := make(chan []dto.CoinTransaction)
	errCh := make(chan *dto.Error, 4)

	go fetchUser(userId, userCh, errCh)
	go fetchInventory(userId, invCh, errCh)
	go fetchTransactions(userId, receivedCh, sentCh, errCh)

	user := <-userCh
	if user == nil {
		return nil, <-errCh
	}
	inventory := <-invCh
	received := <-receivedCh
	sent := <-sentCh

	select {
	case err := <-errCh:
		return nil, err
	default:
	}

	return &dto.InfoResponse{
		Coins:     user.Coins,
		Inventory: inventory,
		CoinHistory: dto.CoinHistory{
			Received: received,
			Sent:     sent,
		},
	}, nil
}

func fetchUser(userId int, ch chan<- *models.Employee, errCh chan<- *dto.Error) {
	var user models.Employee
	err := db.DB.QueryRow("SELECT id, coins FROM employees WHERE id = $1", userId).Scan(&user.ID, &user.Coins)
	if err != nil {
		if err == sql.ErrNoRows {
			errCh <- &dto.Error{Code: "USER_NOT_FOUND", Message: "Пользователь не найден", StatusCode: 404}
		} else {
			errCh <- &dto.Error{Code: "DB_ERROR", Message: "Ошибка базы данных", StatusCode: 500}
		}
		ch <- nil
		return
	}
	ch <- &user
}

func fetchInventory(userId int, ch chan<- []dto.InventoryItem, errCh chan<- *dto.Error) {
	rows, err := db.DB.Query(`SELECT items.name, COUNT(*) FROM purchases JOIN items ON purchases.item_id = items.id WHERE purchases.employee_id = $1 GROUP BY items.name`, userId)
	if err != nil {
		errCh <- &dto.Error{Code: "ITEMS_RETRIEVAL_ERROR", Message: "Ошибка при извлечении инвентаря", StatusCode: 500}
		ch <- nil
		return
	}
	defer rows.Close()

	var inventory []dto.InventoryItem
	for rows.Next() {
		var item dto.InventoryItem
		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			errCh <- &dto.Error{Code: "ITEMS_SCAN_ERROR", Message: "Ошибка при обработке инвентаря", StatusCode: 500}
			ch <- nil
			return
		}
		inventory = append(inventory, item)
	}
	ch <- inventory
}

func fetchTransactions(userId int, receivedCh, sentCh chan<- []dto.CoinTransaction, errCh chan<- *dto.Error) {
	rows, err := db.DB.Query(`SELECT sender_id, receiver_id, amount FROM transactions WHERE sender_id = $1 OR receiver_id = $1`, userId)
	if err != nil {
		errCh <- &dto.Error{Code: "TRANSACTIONS_RETRIEVAL_ERROR", Message: "Ошибка при извлечении транзакций", StatusCode: 500}
		receivedCh <- nil
		sentCh <- nil
		return
	}
	defer rows.Close()

	var received, sent []dto.CoinTransaction
	for rows.Next() {
		var t dto.CoinTransaction
		var senderId, receiverId int
		if err := rows.Scan(&senderId, &receiverId, &t.Amount); err != nil {
			errCh <- &dto.Error{Code: "TRANSACTIONS_SCAN_ERROR", Message: "Ошибка при обработке транзакций", StatusCode: 500}
			receivedCh <- nil
			sentCh <- nil
			return
		}
		if senderId == userId {
			t.ToUser = fetchUsername(receiverId, errCh)
			sent = append(sent, t)
		} else {
			t.FromUser = fetchUsername(senderId, errCh)
			received = append(received, t)
		}
	}
	receivedCh <- received
	sentCh <- sent
}

func fetchUsername(userId int, errCh chan<- *dto.Error) string {
	var username string
	err := db.DB.QueryRow("SELECT username FROM employees WHERE id = $1", userId).Scan(&username)
	if err != nil {
		errCh <- &dto.Error{Code: "USERNAME_FETCH_ERROR", Message: "Ошибка при получении имени пользователя", StatusCode: 500}
		return ""
	}
	return username
}
