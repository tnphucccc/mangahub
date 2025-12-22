package models

import (
	"time"
)

// User represents a user in the system.
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// UserRegisterRequest represents the request body for user registration.
type UserRegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserLoginRequest represents the request body for user login.
type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents the response from successful authentication (login or register).
type AuthResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

// ErrorResponse represents a generic error response from the API.
type ErrorResponse struct {
	Error string `json:"error"`
}
