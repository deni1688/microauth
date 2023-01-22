package main

import (
	"fmt"
	"log"
)

type adminService struct {
	storage     Storage
	authService AuthService
	encryption  Encryption
}

func NewAdminService(s Storage, e Encryption, as AuthService) AdminService {
	return &adminService{storage: s, encryption: e, authService: as}
}

func (s adminService) SaveAdmin(r SaveAdminParams) error {
	a := Admin{
		ID:        r.ID,
		Firstname: r.Firstname,
		Lastname:  r.Lastname,
		Email:     r.Email,
	}

	if r.Password != "" && r.ID != 0 {
		hash, err := s.encryption.Hash(r.Password)
		if err != nil {
			log.Printf("error: hash password failed %v\n", err)
			return fmt.Errorf("admin create failed")
		}

		a.PasswordHash = hash
		if err := s.authService.Expire(a.AuthToken.ID); err != nil {
			log.Printf("error: expire token on admin save failed %v\n", err)
		}
	}

	if err := s.storage.Save(a); err != nil {
		log.Printf("error: save admin %v\n", err)
		return fmt.Errorf("save admin failed")
	}

	return nil
}

func (s adminService) ListAdmins() ([]Admin, error) {
	list, err := s.storage.FindAll()
	if err != nil {
		log.Printf("error: find all admins %v\n", err)
		return list, fmt.Errorf("find all admins failed")
	}

	return list, nil
}

func (s adminService) RemoveAdmin(id AdminID) error {
	if err := s.storage.DeleteByID(id); err != nil {
		log.Printf("error: remove admin %v\n", err)
		return fmt.Errorf("remove admin failed")
	}

	return nil
}
