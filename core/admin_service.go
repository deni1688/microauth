package core

import (
	"context"
	"fmt"
	"log"
)

type adminService struct {
	storage     Storage
	authService AuthService
	hasher      Hasher
}

func NewAdminService(s Storage, h Hasher, as AuthService) AdminService {
	return &adminService{storage: s, hasher: h, authService: as}
}

func (s adminService) SaveAdmin(ctx context.Context, r SaveParams) error {
	a, err := NewAdmin(r)
	if err != nil {
		log.Printf("error: core from save params %v\n", err)
		return fmt.Errorf("core from save params failed")
	}

	if r.Password != "" && a.HashPassword(s.hasher, r.Password) != nil {
		return fmt.Errorf("hash password failed")
	}

	if r.ID != 0 && s.authService.Expire(ctx, a.AuthToken.ID) != nil {
		return fmt.Errorf("expire auth token failed")
	}

	if err = s.storage.Save(ctx, a); err != nil {
		log.Printf("error: save core %v\n", err)
		return fmt.Errorf("save core failed")
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
		log.Printf("error: remove core %v\n", err)
		return fmt.Errorf("remove core failed")
	}

	return nil
}
