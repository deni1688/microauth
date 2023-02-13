package infra

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{}

func (be BcryptHasher) Hash(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(bytes), err
}

func (be BcryptHasher) Compare(pass, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) == nil
}
