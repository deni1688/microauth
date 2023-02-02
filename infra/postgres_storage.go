package infra

import (
	"context"
	"gorm.io/gorm"
	"microauth/core"
)

type postgresStorage struct {
	*gorm.DB
}

type admin struct {
	gorm.Model
	Firstname string
	Lastname  string
	Email     string `gorm:"unique"`
	Password  string
	AuthToken core.AuthToken `gorm:"embedded;embeddedPrefix:auth_token_"`
}

func toAdmin(a *admin) core.Admin {
	return core.Admin{
		ID:           core.AdminID(a.ID),
		Firstname:    a.Firstname,
		Lastname:     a.Lastname,
		Email:        a.Email,
		PasswordHash: a.Password,
		AuthToken:    a.AuthToken,
	}
}

func fromAdmin(a core.Admin) *admin {
	return &admin{
		Firstname: a.Firstname,
		Lastname:  a.Lastname,
		Email:     a.Email,
		Password:  a.PasswordHash,
		AuthToken: a.AuthToken,
	}
}

func NewPostgresStorage(db *gorm.DB) (core.Storage, error) {
	if err := db.AutoMigrate(&admin{}); err != nil {
		return nil, err
	}

	return &postgresStorage{db}, nil
}

func (s postgresStorage) Save(ctx context.Context, a core.Admin) error {
	row := fromAdmin(a)
	if a.ID == core.AdminID(0) {
		return s.Create(row).Error
	}

	return s.WithContext(ctx).Model(&admin{}).Where("id = ?", a.ID).
		Updates(row).
		Error
}

func (s postgresStorage) FindAll(ctx context.Context) ([]core.Admin, error) {
	var rows []admin
	var admins []core.Admin

	tx := s.WithContext(ctx).Find(&rows)
	if tx.Error != nil {
		return admins, tx.Error
	}

	for _, row := range rows {
		admins = append(admins, toAdmin(&row))
	}

	return admins, nil
}

func (s postgresStorage) FindByID(ctx context.Context, id core.AdminID) (core.Admin, error) {
	var row admin

	tx := s.First(ctx, &row, id)
	if tx.Error != nil {
		return toAdmin(&row), tx.Error
	}

	return toAdmin(&row), nil
}

func (s postgresStorage) FindByEmail(ctx context.Context, email string) (core.Admin, error) {
	var row admin

	tx := s.WithContext(ctx).First(&row, "email = ?", email)
	if tx.Error != nil {
		return toAdmin(&row), tx.Error
	}

	return toAdmin(&row), nil
}

func (s postgresStorage) FindByAuthTokenID(ctx context.Context, id core.AuthTokenID) (core.Admin, error) {
	var row admin

	tx := s.WithContext(ctx).First(&row, "auth_token_id = ?", id)
	if tx.Error != nil {
		return toAdmin(&row), tx.Error
	}

	return toAdmin(&row), nil
}

func (s postgresStorage) DeleteByID(ctx context.Context, id core.AdminID) error {
	return s.WithContext(ctx).Delete(&admin{}, id).Error
}