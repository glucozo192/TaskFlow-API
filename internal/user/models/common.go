package models

type LoginRequest struct {
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type LoginResponse struct {
	User  User
	Token string `db:"token" json:"token"`
}
