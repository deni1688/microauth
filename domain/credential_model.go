package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

type Credential struct {
	ID           CredentialID `json:"id"`
	Name         string       `json:"name"`
	PasswordHash string       `json:"-"`
	AuthToken    AuthToken    `json:"-"`
}

func NewCredential(params SaveParams) (Credential, error) {
	if err := validate(params); err != nil {
		return Credential{}, err
	}

	return Credential{
		ID:   params.ID,
		Name: params.Name,
	}, nil
}

func validate(params SaveParams) error {
	if params.Name == "" {
		return fmt.Errorf("name is required")
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

func (a *Credential) HashPassword(hasher Hasher, password string) error {
	hash, err := hasher.Hash(password)
	if err != nil {
		return fmt.Errorf("hash password failed %v", err)
	}

	a.PasswordHash = hash
	return nil
}

func (a *Credential) GenerateAuthToken() error {
	h := sha256.New()
	if _, err := h.Write([]byte(fmt.Sprintf("%d-%s", time.Now().Unix(), randString()))); err != nil {
		return fmt.Errorf("sha256 write for token id failed %v", err)
	}

	a.AuthToken = AuthToken{
		ID:        AuthTokenID(hex.EncodeToString(h.Sum(nil))),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	return nil
}

func (a *Credential) AuthTokenExpired() bool {
	return a.AuthToken.ExpiresAt.Before(time.Now())
}

func (a *Credential) ExpireAuthToken() {
	a.AuthToken = AuthToken{ID: "-", ExpiresAt: time.Now().Add(-(time.Hour * 24))}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randString() string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 12)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
