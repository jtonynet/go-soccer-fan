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

type ExternalApi struct {
	URL   string
	Token string
}

type Config struct {
	Database    *Database
	ExternalApi *ExternalApi
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
		ExternalApi: &ExternalApi{
			URL:   os.Getenv("EXTERNAL_API_URL"),
			Token: os.Getenv("EXTERNAL_API_TOKEN"),
		},
	}
}
