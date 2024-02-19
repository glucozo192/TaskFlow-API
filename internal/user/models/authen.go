package models

import (
	"time"
)

type CreateUserRequest struct {
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	ID                string    `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResponse struct {
	//SessionID             uuid.UUID    `json:"session_id"`
	AccessToken string `json:"access_token"`
	//AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	//RefreshToken          string       `json:"refresh_token"`
	//RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User UserResponse `json:"user"`
}
