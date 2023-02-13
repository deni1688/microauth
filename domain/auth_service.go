package domain

import (
	"context"
	"fmt"
	"log"
)

type authService struct {
	storage Storage
	hasher  Hasher
}

func NewAuthService(s Storage, h Hasher) AuthService {
	return &authService{storage: s, hasher: h}
}

func (s authService) Authenticate(ctx context.Context, r AuthParams) (AuthTokenID, error) {
	a, err := s.storage.FindByName(ctx, r.Name)
	if err != nil {
		log.Printf("error: find credential by name %v\n", err)
		return "", fmt.Errorf("find credential by name failed")
	}

	if !s.hasher.Compare(r.Password, a.PasswordHash) {
		log.Println("error: invalid password")
		return "", fmt.Errorf("invalid password")
	}

	if err = a.GenerateAuthToken(); err != nil {
		log.Printf("error: generate auth token %v\n", err)
		return "", fmt.Errorf("generate auth token failed")
	}

	if err = s.storage.Save(ctx, a); err != nil {
		log.Printf("error: save credential %v\n", err)
		return "", fmt.Errorf("save credential failed")
	}

	return a.AuthToken.ID, nil
}

func (s authService) Validate(ctx context.Context, id AuthTokenID) error {
	if id == "-" {
		return fmt.Errorf("token invalid")
	}

	a, err := s.storage.FindByAuthTokenID(ctx, id)
	if err != nil {
		log.Printf("error: find credential by token id %v\n", err)
		return fmt.Errorf("find credential by token id failed")
	}

	if a.AuthTokenExpired() {
		return fmt.Errorf("token expired")
	}

	return nil
}

func (s authService) Expire(ctx context.Context, id AuthTokenID) error {
	a, err := s.storage.FindByAuthTokenID(ctx, id)
	if err != nil {
		log.Printf("error: find credential by token id %v\n", err)
		return fmt.Errorf("find credential by token id failed")
	}

	a.ExpireAuthToken()
	if err = s.storage.Save(ctx, a); err != nil {
		log.Printf("error: saving credential on token expire %v\n", err)
		return fmt.Errorf("saving credential on token expire failed")
	}

	return nil
}
