package main

type AdminService interface {
	SaveAdmin(SaveAdminParams) error
	ListAdmins() ([]Admin, error)
	RemoveAdmin(AdminID) error
}
type AuthService interface {
	Authenticate(AuthParams) (AuthTokenID, error)
	Validate(AuthTokenID) error
	Expire(AuthTokenID) error
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
	Hash(string) (string, error)
	Compare(password, hash string) bool
}
