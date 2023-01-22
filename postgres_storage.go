package main

import (
	"gorm.io/gorm"
)

type postgresStorage struct {
	*gorm.DB
}

type admin struct {
	gorm.Model
	Firstname string
	Lastname  string
	Email     string
	Password  string
	AuthToken AuthToken `gorm:"embedded;embeddedPrefix:auth_token_"`
}

func toAdmin(a *admin) Admin {
	return Admin{
		ID:           AdminID(a.ID),
		Firstname:    a.Firstname,
		Lastname:     a.Lastname,
		Email:        a.Email,
		PasswordHash: a.Password,
		AuthToken:    a.AuthToken,
	}
}

func fromAdmin(a Admin) *admin {
	return &admin{
		Firstname: a.Firstname,
		Lastname:  a.Lastname,
		Email:     a.Email,
		Password:  a.PasswordHash,
		AuthToken: a.AuthToken,
	}
}

func NewPostgresStorage(db *gorm.DB) (Storage, error) {
	if err := db.AutoMigrate(&admin{}); err != nil {
		return nil, err
	}

	return &postgresStorage{db}, nil
}

func (s postgresStorage) Save(a Admin) error {
	row := fromAdmin(a)
	if a.ID == AdminID(0) {
		return s.Create(row).Error
	}

	return s.Model(&admin{}).Where("id = ?", a.ID).
		Updates(row).
		Error
}

func (s postgresStorage) FindAll() ([]Admin, error) {
	var rows []admin
	var admins []Admin

	tx := s.Find(&rows)
	if tx.Error != nil {
		return admins, tx.Error
	}

	for _, row := range rows {
		admins = append(admins, toAdmin(&row))
	}

	return admins, nil
}

func (s postgresStorage) FindByID(id AdminID) (Admin, error) {
	var row admin

	tx := s.First(&row, id)
	if tx.Error != nil {
		return toAdmin(&row), tx.Error
	}

	return toAdmin(&row), nil
}

func (s postgresStorage) FindByEmail(email string) (Admin, error) {
	var row admin

	tx := s.First(&row, "email = ?", email)
	if tx.Error != nil {
		return toAdmin(&row), tx.Error
	}

	return toAdmin(&row), nil
}

func (s postgresStorage) FindByAuthTokenID(id AuthTokenID) (Admin, error) {
	var row admin

	tx := s.First(&row, "auth_token_id = ?", id)
	if tx.Error != nil {
		return toAdmin(&row), tx.Error
	}

	return toAdmin(&row), nil
}

func (s postgresStorage) DeleteByID(id AdminID) error {
	return s.Delete(&admin{}, id).Error
}
