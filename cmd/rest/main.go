package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"microauth/core"
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

	authSrv := core.NewAuthService(ps, bh)
	adminSrv := core.NewAdminService(ps, bh, authSrv)

	mustCreateDefaultAdmin(cf, ps, adminSrv)

	adminH := rest.NewAdminHandler(adminSrv)
	authH := rest.NewAuthHandler(authSrv)
	authM := rest.NewAuthMiddleware(authSrv)

	e := echo.New()

	api := e.Group("/api/v1")

	api.POST("/login", authH.HandleLogin)
	api.POST("/logout", authH.HandleLogout)

	dash := api.Group("/dashboard")
	dash.Use(authM)

	dash.GET("/admins", adminH.HandleGetAdmins)
	dash.POST("/admins", adminH.HandleSaveAdmin)
	dash.DELETE("/admins/:id", adminH.HandleDeleteAdmin)

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

func mustCreateDefaultAdmin(cf *config, s core.Storage, as core.AdminService) {
	ctx := context.Background()
	_, err := s.FindByEmail(ctx, cf.DefaultAdminEmail)
	if err == nil {
		log.Println("Skipping default core creation")
		return
	}

	if err = as.SaveAdmin(ctx, core.SaveParams{
		Firstname: cf.DefaultAdminFirstname,
		Lastname:  cf.DefaultAdminLastname,
		Email:     cf.DefaultAdminEmail,
		Password:  cf.DefaultAdminPassword,
	}); err != nil {
		log.Fatalf("Failed to create default core: %v", err)
	}
}
