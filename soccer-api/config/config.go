package config

import (
	"log"
	"os"
	"strconv"
)

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

type MailNotification struct {
	SMTPHost         string
	SMTPPort         int
	SMTPUser         string
	SMTPPassword     string
	EmailFromEmail   string
	EmailFromName    string
	UseMailTLS       bool
	IsDevelopmentEnv bool
}

type Config struct {
	Database         *Database
	ExternalApi      *ExternalApi
	MailNotification *MailNotification
}

func LoadConfig() *Config {
	smtpPort, err := strconv.Atoi(os.Getenv("MAIL_NOTIFICATION_SMTP_PORT"))
	if err != nil {
		log.Fatalf("Erro ao converter MAIL_NOTIFICATION_SMTP_PORT para inteiro: %v", err)
	}

	useMailTLS, err := strconv.ParseBool(os.Getenv("MAIL_NOTIFICATION_MAIL_TLS"))
	if err != nil {
		log.Fatalf("Erro ao converter MAIL_NOTIFICATION_MAIL_TLS para booleano: %v", err)
	}

	isDevelopmentEnv, err := strconv.ParseBool(os.Getenv("MAIL_NOTIFICATION_IS_DEVELOPMENT_ENV"))
	if err != nil {
		log.Fatalf("Erro ao converter MAIL_NOTIFICATION_IS_DEVELOPMENT_ENV para booleano: %v", err)
	}

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
		MailNotification: &MailNotification{
			SMTPHost:         os.Getenv("MAIL_NOTIFICATION_SMTP_HOST"),
			SMTPPort:         smtpPort,
			SMTPUser:         os.Getenv("MAIL_NOTIFICATION_SMTP_USER"),
			SMTPPassword:     os.Getenv("MAIL_NOTIFICATION_SMTP_PASSWORD"),
			EmailFromEmail:   os.Getenv("MAIL_NOTIFICATION_EMAILS_FROM_EMAIL"),
			EmailFromName:    os.Getenv("MAIL_NOTIFICATION_EMAILS_FROM_NAME"),
			UseMailTLS:       useMailTLS,
			IsDevelopmentEnv: isDevelopmentEnv,
		},
	}
}
