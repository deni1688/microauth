package main

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptEncryption struct{}

func (be BcryptEncryption) Hash(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(bytes), err
}

func (be BcryptEncryption) Compare(pass, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) == nil
}
