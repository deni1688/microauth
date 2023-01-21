package main

import "time"

type AdminID uint
type AuthTokenID string

type AdminBase struct {
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type AuthToken struct {
	ID        AuthTokenID
	ExpiresAt time.Time
}
type Admin struct {
	AdminBase
	ID        AdminID   `json:"id"`
	AuthToken AuthToken `json:"-"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SaveRequest struct {
	ID        AdminID `json:"id"`
	Firstname string  `json:"firstname"`
	Lastname  string  `json:"lastname"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
}
