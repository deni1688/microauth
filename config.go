package main

import "os"

type config struct {
	PostgresHost          string
	PostgresUser          string
	PostgresPassword      string
	PostgresDBName        string
	PostgresPort          string
	PostgresSSLMode       string
	PostgresTimezone      string
	DefaultAdminFirstname string
	DefaultAdminLastname  string
	DefaultAdminEmail     string
	DefaultAdminPassword  string
}

func loadConfig() *config {
	return &config{
		PostgresHost:          os.Getenv("DB_HOST"),
		PostgresUser:          os.Getenv("DB_USER"),
		PostgresPassword:      os.Getenv("DB_PASS"),
		PostgresDBName:        os.Getenv("DB_NAME"),
		PostgresPort:          os.Getenv("DB_PORT"),
		PostgresSSLMode:       os.Getenv("DB_SSL_MODE"),
		PostgresTimezone:      os.Getenv("DB_TIMEZONE"),
		DefaultAdminFirstname: os.Getenv("ADMIN_FIRSTNAME"),
		DefaultAdminLastname:  os.Getenv("ADMIN_LASTNAME"),
		DefaultAdminEmail:     os.Getenv("ADMIN_EMAIL"),
		DefaultAdminPassword:  os.Getenv("ADMIN_PASSWORD"),
	}
}
