package main

import (
	"context"
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

func (s adminService) SaveAdmin(ctx context.Context, r SaveParams) error {
	a, err := NewAdmin(r)
	if err != nil {
		log.Printf("error: admin from save params %v\n", err)
		return fmt.Errorf("admin from save params failed")
	}

	if r.Password != "" && r.ID != 0 {
		if err = a.HashPassword(s.encryption, r.Password); err != nil {
			log.Printf("error: hash password failed %v\n", err)
			return fmt.Errorf("admin create failed")
		}

		if err = s.authService.Expire(ctx, a.AuthToken.ID); err != nil {
			log.Printf("error: expire token on admin save failed %v\n", err)
		}
	}

	if err = s.storage.Save(ctx, a); err != nil {
		log.Printf("error: save admin %v\n", err)
		return fmt.Errorf("save admin failed")
	}

	return nil
}

func (s adminService) ListAdmins(ctx context.Context) ([]Admin, error) {
	list, err := s.storage.FindAll(ctx)
	if err != nil {
		log.Printf("error: find all admins %v\n", err)
		return list, fmt.Errorf("find all admins failed")
	}

	return list, nil
}

func (s adminService) RemoveAdmin(ctx context.Context, id AdminID) error {
	if err := s.storage.DeleteByID(ctx, id); err != nil {
		log.Printf("error: remove admin %v\n", err)
		return fmt.Errorf("remove admin failed")
	}

	return nil
}
