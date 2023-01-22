package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	db := mustConnectDB()
	ps := MustNewPostgresStorage(db)
	bc := BcryptEncryption{}

	authSrv := NewAuthService(ps, bc)
	adminSrv := NewAdminService(ps, bc, authSrv)

	_ = adminSrv.SaveAdmin(SaveAdminParams{
		Firstname: "admin",
		Lastname:  "admin",
		Email:     "admin@test.com",
		Password:  "admin",
	})

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

func mustConnectDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(fmt.Sprintf(`
	host=%s
	user=%s
	password=%s
	dbname=%s
	port=%s
	sslmode=%s
	TimeZone=%s`,
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL_MODE"),
		os.Getenv("DB_TIMEZONE"),
	)), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}
