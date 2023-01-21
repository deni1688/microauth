package main

type AdminService interface {
	SaveAdmin(SaveRequest) error
	ListAdmins() ([]Admin, error)
	RemoveAdmin(AdminID) error
}
type AuthService interface {
	Authenticate(AuthRequest) (AuthTokenID, error)
	Validate(AuthTokenID) error
	Invalidate(AuthTokenID) error
}

type Storage interface {
	Save(Admin) error
	FindAll() ([]Admin, error)
	FindByID(AdminID) (Admin, error)
	FindByEmail(string) (Admin, error)
	FindByAuthTokenID(AuthTokenID) (Admin, error)
	DeleteByID(AdminID) error
}

type Encryption interface {
	Hash(password string) (string, error)
	Compare(password string, hash string) bool
}
