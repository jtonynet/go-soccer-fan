package service

import (
	"crypto/tls"
	"log"

	"github.com/jtonynet/go-soccer-fan/soccer-api/config"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
	UseTLS   bool
}

var isDevelopmentEnv = true

func NewEmailService(cfg *config.MailNotification) *EmailService {

	isDevelopmentEnv = cfg.IsDevelopmentEnv

	return &EmailService{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Username: cfg.SMTPUser,
		Password: cfg.SMTPPassword,
		From:     cfg.EmailFromEmail,
		FromName: cfg.EmailFromName,
		UseTLS:   cfg.UseMailTLS,
	}
}

func (es *EmailService) SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", es.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(es.Host, es.Port, es.Username, es.Password)

	// Disable certificate verification for development purposes
	if es.UseTLS && isDevelopmentEnv {
		d.TLSConfig = &tls.Config{
			InsecureSkipVerify: true, // Skip certificate verification
		}
	} else {
		d.TLSConfig = nil
	}

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Println("Email sent successfully")
	return nil
}
