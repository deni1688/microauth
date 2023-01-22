package main

import (
	"fmt"
	"time"
)

type Admin struct {
	ID           AdminID   `json:"id"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	AuthToken    AuthToken `json:"-"`
}

func NewAdmin(params SaveParams) (Admin, error) {
	if params.Email == "" {
		return Admin{}, fmt.Errorf("email is required")
	}

	if params.Firstname == "" {
		return Admin{}, fmt.Errorf("firstname is required")
	}

	if params.Lastname == "" {
		return Admin{}, fmt.Errorf("lastname is required")
	}

	if params.ID == 0 {
		if params.Password == "" {
			return Admin{}, fmt.Errorf("password is required")
		}
	}

	if params.Password != "" && len(params.Password) < 8 {
		return Admin{}, fmt.Errorf("password must be at least 8 characters")
	}

	return Admin{
		ID:        params.ID,
		Firstname: params.Firstname,
		Lastname:  params.Lastname,
		Email:     params.Email,
	}, nil
}

func (a *Admin) HashPassword(encryption Encryption, password string) error {
	hash, err := encryption.Hash(password)
	if err != nil {
		return fmt.Errorf("hash password failed")
	}

	a.PasswordHash = hash
	return nil
}

func (a *Admin) ExpireAuthToken() {
	a.AuthToken = AuthToken{ID: "-", ExpiresAt: time.Now().Add(-(time.Hour * 24))}
}
