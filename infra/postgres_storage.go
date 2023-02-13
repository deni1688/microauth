package infra

import (
	"context"
	"gorm.io/gorm"
	"microauth/domain"
)

type postgresStorage struct {
	*gorm.DB
}

type credential struct {
	gorm.Model
	Name      string `gorm:"unique"`
	Password  string
	AuthToken domain.AuthToken `gorm:"embedded;embeddedPrefix:auth_token_"`
}

func toCredential(a *credential) domain.Credential {
	return domain.Credential{
		ID:           domain.CredentialID(a.ID),
		Name:         a.Name,
		PasswordHash: a.Password,
		AuthToken:    a.AuthToken,
	}
}

func fromCredential(a domain.Credential) *credential {
	return &credential{
		Name:      a.Name,
		Password:  a.PasswordHash,
		AuthToken: a.AuthToken,
	}
}

func NewPostgresStorage(db *gorm.DB) (domain.Storage, error) {
	if err := db.AutoMigrate(&credential{}); err != nil {
		return nil, err
	}

	return &postgresStorage{db}, nil
}

func (s postgresStorage) Save(ctx context.Context, a domain.Credential) error {
	row := fromCredential(a)
	if a.ID == domain.CredentialID(0) {
		return s.Create(row).Error
	}

	return s.WithContext(ctx).Model(&credential{}).Where("id = ?", a.ID).
		Updates(row).
		Error
}

func (s postgresStorage) FindAll(ctx context.Context) ([]domain.Credential, error) {
	var rows []credential
	var credentials []domain.Credential

	tx := s.WithContext(ctx).Find(&rows)
	if tx.Error != nil {
		return credentials, tx.Error
	}

	for _, row := range rows {
		credentials = append(credentials, toCredential(&row))
	}

	return credentials, nil
}

func (s postgresStorage) FindByID(ctx context.Context, id domain.CredentialID) (domain.Credential, error) {
	var row credential

	tx := s.First(ctx, &row, id)
	if tx.Error != nil {
		return toCredential(&row), tx.Error
	}

	return toCredential(&row), nil
}

func (s postgresStorage) FindByName(ctx context.Context, name string) (domain.Credential, error) {
	var row credential

	tx := s.WithContext(ctx).First(&row, "name = ?", name)
	if tx.Error != nil {
		return toCredential(&row), tx.Error
	}

	return toCredential(&row), nil
}

func (s postgresStorage) FindByAuthTokenID(ctx context.Context, id domain.AuthTokenID) (domain.Credential, error) {
	var row credential

	tx := s.WithContext(ctx).First(&row, "auth_token_id = ?", id)
	if tx.Error != nil {
		return toCredential(&row), tx.Error
	}

	return toCredential(&row), nil
}

func (s postgresStorage) DeleteByID(ctx context.Context, id domain.CredentialID) error {
	return s.WithContext(ctx).Delete(&credential{}, id).Error
}
