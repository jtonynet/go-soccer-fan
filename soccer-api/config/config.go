package config

import (
	"log"
	"os"
	"strconv"
)

type API struct {
	Port             int
	SecretKey        string
	TokenLifespanInH int
}

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

type RabbitMQ struct {
	Host                  string
	User                  string
	Password              string
	Port                  string
	MatchNoticationsQueue string
	FanNotificationsQueue string
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
	API              *API
	Database         *Database
	ExternalApi      *ExternalApi
	RabbitMQ         *RabbitMQ
	MailNotification *MailNotification
}

func LoadConfig() *Config {
	apiPort, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		log.Fatalf("Erro ao converter API_PORT para inteiro: %v", err)
	}

	apiTokenLifespanInH, err := strconv.Atoi(os.Getenv("API_TOKEN_LIFESPAN_IN_H"))
	if err != nil {
		log.Fatalf("Erro ao converter API_TOKEN_LIFESPAN_IN_H para inteiro: %v", err)
	}

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
		API: &API{
			Port:             apiPort,
			SecretKey:        os.Getenv("API_SECRET_KEY"),
			TokenLifespanInH: apiTokenLifespanInH,
		},
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
		RabbitMQ: &RabbitMQ{
			User:                  os.Getenv("RABBITMQ_USER"),
			Password:              os.Getenv("RABBITMQ_PASS"),
			Port:                  os.Getenv("RABBITMQ_PORT"),
			Host:                  os.Getenv("RABBITMQ_HOST"),
			MatchNoticationsQueue: os.Getenv("RABBITMQ_QUEUE_MATCH_NOTIFICATIONS"),
			FanNotificationsQueue: os.Getenv("RABBITMQ_QUEUE_FAN_MATCH_NOTIFICATIONS"),
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
