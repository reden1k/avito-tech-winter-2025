package dto

type SendCoinsRequest struct {
	ReceiverUsername string `json:"toUser"`
	Amount           int    `json:"amount"`
}

type SendCoinsResponse struct {
	Message string `json:"message"`
}
