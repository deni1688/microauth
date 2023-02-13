package main

import "os"

type config struct {
	PostgresHost              string
	PostgresUser              string
	PostgresPassword          string
	PostgresDBName            string
	PostgresPort              string
	PostgresSSLMode           string
	PostgresTimezone          string
	DefaultCredentialName     string
	DefaultCredentialPassword string
}

func loadConfig() *config {
	return &config{
		PostgresHost:              os.Getenv("DB_HOST"),
		PostgresUser:              os.Getenv("DB_USER"),
		PostgresPassword:          os.Getenv("DB_PASS"),
		PostgresDBName:            os.Getenv("DB_NAME"),
		PostgresPort:              os.Getenv("DB_PORT"),
		PostgresSSLMode:           os.Getenv("DB_SSL_MODE"),
		PostgresTimezone:          os.Getenv("DB_TIMEZONE"),
		DefaultCredentialName:     os.Getenv("CREDENTIAL_NAME"),
		DefaultCredentialPassword: os.Getenv("CREDENTIAL_PASSWORD"),
	}
}
