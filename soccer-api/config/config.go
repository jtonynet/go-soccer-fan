package config

import "os"

type Database struct {
	Host     string
	User     string
	Password string
	DBname   string
	Port     string
	SSLmode  string
}

type Config struct {
	Database *Database
}

func LoadConfig() *Config {
	return &Config{
		Database: &Database{
			Host:     os.Getenv("DATABASE_HOST"),
			User:     os.Getenv("DATABASE_USER"),
			Password: os.Getenv("DATABASE_PASSWORD"),
			DBname:   os.Getenv("DATABASE_DB"),
			Port:     os.Getenv("DATABASE_PORT"),
			SSLmode:  os.Getenv("DATABASE_SSLMODE"),
		},
	}
}
