package main

import (
	"context"
	"fmt"
	"log"
)

type authService struct {
	storage    Storage
	encryption Encryption
}

func NewAuthService(s Storage, e Encryption) AuthService {
	return &authService{storage: s, encryption: e}
}

func (s authService) Authenticate(ctx context.Context, r AuthParams) (AuthTokenID, error) {
	a, err := s.storage.FindByEmail(ctx, r.Email)
	if err != nil {
		log.Printf("error: find admin by email %v\n", err)
		return "", fmt.Errorf("find admin by email failed")
	}

	if !s.encryption.Compare(r.Password, a.PasswordHash) {
		log.Println("error: invalid password")
		return "", fmt.Errorf("invalid password")
	}

	if err = a.GenerateAuthToken(); err != nil {
		log.Printf("error: generate auth token %v\n", err)
		return "", fmt.Errorf("generate auth token failed")
	}

	if err = s.storage.Save(ctx, a); err != nil {
		log.Printf("error: save admin %v\n", err)
		return "", fmt.Errorf("save admin failed")
	}

	return a.AuthToken.ID, nil
}

func (s authService) Validate(ctx context.Context, id AuthTokenID) error {
	if id == "-" {
		return fmt.Errorf("token invalid")
	}

	a, err := s.storage.FindByAuthTokenID(ctx, id)
	if err != nil {
		log.Printf("error: find admin by token id %v\n", err)
		return fmt.Errorf("find admin by token id failed")
	}

	if a.AuthTokenExpired() {
		return fmt.Errorf("token expired")
	}

	return nil
}

func (s authService) Expire(ctx context.Context, id AuthTokenID) error {
	a, err := s.storage.FindByAuthTokenID(ctx, id)
	if err != nil {
		log.Printf("error: find admin by token id %v\n", err)
		return fmt.Errorf("find admin by token id failed")
	}

	a.ExpireAuthToken()
	if err = s.storage.Save(ctx, a); err != nil {
		log.Printf("error: aving admin on token expire %v\n", err)
		return fmt.Errorf("saving admin on token expire failed")
	}

	return nil
}
