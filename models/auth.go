package models

type AuthRequest struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
