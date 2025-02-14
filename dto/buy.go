package dto

type BuyRequest struct {
	ItemName string `json:"item_name"`
}

type BuyResponse struct {
	Message string `json:"message"`
}
