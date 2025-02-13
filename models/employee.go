package models

type Employee struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	PasswordHash string `json:"-"`
	Coins        int    `json:"coins"`
}
