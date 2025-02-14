package models

type Purchase struct {
	ID         int `json:"id"`
	EmployeeID int `json:"employee_id"`
	ItemID     int `json:"item_id"`
}
