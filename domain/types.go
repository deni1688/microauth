package domain

import "time"

type CredentialID uint
type AuthTokenID string
type AuthToken struct {
	ID        AuthTokenID
	ExpiresAt time.Time
}
type AuthParams struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
type SaveParams struct {
	ID       CredentialID `json:"id"`
	Name     string       `json:"name"`
	Password string       `json:"password"`
}
