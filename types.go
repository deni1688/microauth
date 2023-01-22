package main

import "time"

type AdminID uint
type AuthTokenID string

type AuthToken struct {
	ID        AuthTokenID
	ExpiresAt time.Time
}
type Admin struct {
	ID           AdminID   `json:"id"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	AuthToken    AuthToken `json:"-"`
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SaveAdminParams struct {
	ID        AdminID `json:"id"`
	Firstname string  `json:"firstname"`
	Lastname  string  `json:"lastname"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
}
