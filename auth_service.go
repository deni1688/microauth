package main

import (
	"fmt"
	"log"
	"time"
)

type authService struct {
	storage    Storage
	encryption Encryption
}

func NewAuthService(s Storage, e Encryption) AuthService {
	return &authService{storage: s, encryption: e}
}

func (s authService) Authenticate(r AuthParams) (AuthTokenID, error) {
	a, err := s.storage.FindByEmail(r.Email)
	if err != nil {
		log.Printf("error: find admin by email %v\n", err)
		return "", fmt.Errorf("find admin by email failed")
	}

	if !s.encryption.Compare(r.Password, a.PasswordHash) {
		log.Println("error: invalid password")
		return "", fmt.Errorf("invalid password")
	}

	id, err := s.encryption.Hash(fmt.Sprintf("%d-%d", a.ID, time.Now().Unix()))
	if err != nil {
		log.Printf("error: hash token id %v\n", err)
		return "", fmt.Errorf("hash token id failed")
	}

	a.AuthToken = AuthToken{
		ID:        AuthTokenID(id),
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}

	if err = s.storage.Save(a); err != nil {
		log.Printf("error: save admin %v\n", err)
		return "", fmt.Errorf("save admin failed")
	}

	return a.AuthToken.ID, nil
}

func (s authService) Validate(id AuthTokenID) error {
	if id == "-" {
		return fmt.Errorf("token invalid")
	}

	a, err := s.storage.FindByAuthTokenID(id)
	if err != nil {
		log.Printf("error: find admin by token id %v\n", err)
		return fmt.Errorf("find admin by token id failed")
	}

	if a.AuthToken.ExpiresAt.Before(time.Now()) {
		_ = s.Expire(id)
		return fmt.Errorf("token expired")
	}

	return nil
}

func (s authService) Expire(id AuthTokenID) error {
	a, err := s.storage.FindByAuthTokenID(id)
	if err != nil {
		log.Printf("error: find admin by token id %v\n", err)
		return fmt.Errorf("find admin by token id failed")
	}

	a.ExpireAuthToken()
	if err = s.storage.Save(a); err != nil {
		log.Printf("error: aving admin on token expire %v\n", err)
		return fmt.Errorf("saving admin on token expire failed")
	}

	return nil
}
