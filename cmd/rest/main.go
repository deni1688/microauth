package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"microauth/domain"
	"microauth/infra"
	"microauth/rest"
)

func main() {
	cf := loadConfig()
	db := mustConnectDB(cf)
	ps, err := infra.NewPostgresStorage(db)
	if err != nil {
		log.Fatalf("Failed to create postgres: %v", err)
	}
	bh := infra.BcryptHasher{}

	authSrv := domain.NewAuthService(ps, bh)
	credentialSrv := domain.NewCredentialService(ps, bh, authSrv)

	mustCreateDefaultCredential(cf, ps, credentialSrv)

	credentialH := rest.NewCredentialHandler(credentialSrv)
	authH := rest.NewAuthHandler(authSrv)
	authM := rest.NewAuthMiddleware(authSrv)

	e := echo.New()

	api := e.Group("/api/v1")

	api.POST("/login", authH.HandleLogin)
	api.POST("/logout", authH.HandleLogout)

	dash := api.Group("/dashboard")
	dash.Use(authM)

	dash.GET("/credentials", credentialH.HandleGetCredentials)
	dash.POST("/credentials", credentialH.HandleSaveCredential)
	dash.DELETE("/credentials/:id", credentialH.HandleDeleteCredential)

	e.Logger.Fatal(e.Start(":9876"))
}

func mustConnectDB(cf *config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(fmt.Sprintf(`
	host=%s
	user=%s
	password=%s
	dbname=%s
	port=%s
	sslmode=%s
	TimeZone=%s`,
		cf.PostgresHost,
		cf.PostgresUser,
		cf.PostgresPassword,
		cf.PostgresDBName,
		cf.PostgresPort,
		cf.PostgresSSLMode,
		cf.PostgresTimezone,
	)), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}

func mustCreateDefaultCredential(cf *config, s domain.Storage, as domain.CredentialService) {
	ctx := context.Background()
	_, err := s.FindByName(ctx, cf.DefaultCredentialName)
	if err == nil {
		log.Println("Skipping default core creation")
		return
	}

	if err = as.SaveCredential(ctx, domain.SaveParams{
		Name:     cf.DefaultCredentialName,
		Password: cf.DefaultCredentialPassword,
	}); err != nil {
		log.Fatalf("Failed to create default credential: %v", err)
	} else {
		log.Println("Default credential created")
	}
}
