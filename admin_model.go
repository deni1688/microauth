package main

import (
	"crypto/sha256"
	"encoding/hex"
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
	if err := validate(params); err != nil {
		return Admin{}, err
	}

	return Admin{
		ID:        params.ID,
		Firstname: params.Firstname,
		Lastname:  params.Lastname,
		Email:     params.Email,
	}, nil
}

func validate(params SaveParams) error {
	if params.Email == "" {
		return fmt.Errorf("email is required")
	}

	if params.Firstname == "" {
		return fmt.Errorf("firstname is required")
	}

	if params.Lastname == "" {
		return fmt.Errorf("lastname is required")
	}

	if params.ID == 0 {
		if params.Password == "" {
			return fmt.Errorf("password is required")
		}
	}

	if params.Password != "" && len(params.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	return nil
}

func (a *Admin) HashPassword(encryption Hasher, password string) error {
	hash, err := encryption.Hash(password)
	if err != nil {
		return fmt.Errorf("hash password failed %v", err)
	}

	a.PasswordHash = hash
	return nil
}

func (a *Admin) GenerateAuthToken() error {
	h := sha256.New()
	if _, err := h.Write([]byte(fmt.Sprintf("%d-%d", a.ID, time.Now().Unix()))); err != nil {
		return fmt.Errorf("sha256 write for token id failed %v", err)
	}

	a.AuthToken = AuthToken{
		ID:        AuthTokenID(hex.EncodeToString(h.Sum(nil))),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	return nil
}

func (a *Admin) AuthTokenExpired() bool {
	return a.AuthToken.ExpiresAt.Before(time.Now())
}

func (a *Admin) ExpireAuthToken() {
	a.AuthToken = AuthToken{ID: "-", ExpiresAt: time.Now().Add(-(time.Hour * 24))}
}
