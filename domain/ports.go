package domain

import (
	"context"
)

type CredentialService interface {
	SaveCredential(context.Context, SaveParams) error
	ListCredentials(context.Context) ([]Credential, error)
	RemoveCredential(context.Context, CredentialID) error
}
type AuthService interface {
	Authenticate(context.Context, AuthParams) (AuthTokenID, error)
	Validate(context.Context, AuthTokenID) error
	Expire(context.Context, AuthTokenID) error
}

type Storage interface {
	Save(context.Context, Credential) error
	FindAll(context.Context) ([]Credential, error)
	FindByID(context.Context, CredentialID) (Credential, error)
	FindByName(context.Context, string) (Credential, error)
	FindByAuthTokenID(context.Context, AuthTokenID) (Credential, error)
	DeleteByID(context.Context, CredentialID) error
}

type Hasher interface {
	Hash(string) (string, error)
	Compare(password, hash string) bool
}
