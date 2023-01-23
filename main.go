package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	cf := loadConfig()
	db := mustConnectDB(cf)
	ps, err := NewPostgresStorage(db)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}
	bh := BcryptHasher{}

	authSrv := NewAuthService(ps, bh)
	adminSrv := NewAdminService(ps, bh, authSrv)

	mustCreateDefaultAdmin(cf, ps, adminSrv)

	adminH := NewAdminHandler(adminSrv)
	authH := NewAuthHandler(authSrv)
	authM := NewAuthMiddleware(authSrv)

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

func mustCreateDefaultAdmin(cf *config, s Storage, as AdminService) {
	ctx := context.Background()
	_, err := s.FindByEmail(ctx, cf.DefaultAdminEmail)
	if err == nil {
		log.Println("Skipping default admin creation")
		return
	}

	if err = as.SaveAdmin(ctx, SaveParams{
		Firstname: cf.DefaultAdminFirstname,
		Lastname:  cf.DefaultAdminLastname,
		Email:     cf.DefaultAdminEmail,
		Password:  cf.DefaultAdminPassword,
	}); err != nil {
		log.Fatalf("Failed to create default admin: %v", err)
	}
}
