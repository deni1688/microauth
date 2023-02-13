package domain

import (
	"context"
	"fmt"
	"log"
)

type credentialService struct {
	storage     Storage
	authService AuthService
	hasher      Hasher
}

func NewCredentialService(s Storage, h Hasher, as AuthService) CredentialService {
	return &credentialService{storage: s, hasher: h, authService: as}
}

func (s credentialService) SaveCredential(ctx context.Context, r SaveParams) error {
	a, err := NewCredential(r)
	if err != nil {
		log.Printf("error: core from save params %v\n", err)
		return fmt.Errorf("credential from save params failed")
	}

	if r.Password != "" && a.HashPassword(s.hasher, r.Password) != nil {
		return fmt.Errorf("hash password failed")
	}

	if r.ID != 0 && s.authService.Expire(ctx, a.AuthToken.ID) != nil {
		return fmt.Errorf("expire auth token failed")
	}

	if err = s.storage.Save(ctx, a); err != nil {
		log.Printf("error: save core %v\n", err)
		return fmt.Errorf("save credential failed")
	}

	return nil
}

func (s credentialService) ListCredentials(ctx context.Context) ([]Credential, error) {
	list, err := s.storage.FindAll(ctx)
	if err != nil {
		log.Printf("error: find all credentials %v\n", err)
		return list, fmt.Errorf("find all credentials failed")
	}

	return list, nil
}

func (s credentialService) RemoveCredential(ctx context.Context, id CredentialID) error {
	if err := s.storage.DeleteByID(ctx, id); err != nil {
		log.Printf("error: remove core %v\n", err)
		return fmt.Errorf("remove credential failed")
	}

	return nil
}
