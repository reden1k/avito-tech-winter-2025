package models

type Transaction struct {
	SenderID   int `json:"sender_id"`
	ReceiverID int `json:"receiver_id"`
	Amount     int `json:"amount"`
}
