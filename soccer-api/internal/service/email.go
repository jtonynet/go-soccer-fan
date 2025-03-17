package service

import (
	"crypto/tls"
	"log"

	"github.com/jtonynet/go-soccer-fan/soccer-api/config"
	"gopkg.in/gomail.v2"
)

type Mail struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
	UseTLS   bool
}

var isDevelopmentEnv = true

func NewMail(
	cfg *config.MailNotification,
) *Mail {

	isDevelopmentEnv = cfg.IsDevelopmentEnv

	return &Mail{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Username: cfg.SMTPUser,
		Password: cfg.SMTPPassword,
		From:     cfg.EmailFromEmail,
		FromName: cfg.EmailFromName,
		UseTLS:   cfg.UseMailTLS,
	}
}

func (es *Mail) SendMail(to string, subject string, body string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", es.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(es.Host, es.Port, es.Username, es.Password)

	if es.UseTLS && isDevelopmentEnv {
		d.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	} else {
		d.TLSConfig = nil
	}

	if err := d.DialAndSend(m); err != nil {
		log.Printf("failed to send email: %v", err)
		return err
	}

	log.Println("email sent successfully")
	return nil
}
