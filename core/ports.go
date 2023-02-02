package core

import (
	"context"
)

type AdminService interface {
	SaveAdmin(context.Context, SaveParams) error
	ListAdmins(context.Context) ([]Admin, error)
	RemoveAdmin(context.Context, AdminID) error
}
type AuthService interface {
	Authenticate(context.Context, AuthParams) (AuthTokenID, error)
	Validate(context.Context, AuthTokenID) error
	Expire(context.Context, AuthTokenID) error
}

type Storage interface {
	Save(context.Context, Admin) error
	FindAll(context.Context) ([]Admin, error)
	FindByID(context.Context, AdminID) (Admin, error)
	FindByEmail(context.Context, string) (Admin, error)
	FindByAuthTokenID(context.Context, AuthTokenID) (Admin, error)
	DeleteByID(context.Context, AdminID) error
}

type Hasher interface {
	Hash(string) (string, error)
	Compare(password, hash string) bool
}
