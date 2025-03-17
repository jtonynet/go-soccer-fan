package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/jtonynet/go-soccer-fan/soccer-api/config"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/dto"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/repository"
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

	teamRepo repository.Team
}

var isDevelopmentEnv = true

func NewEmailBroadcast(
	cfg *config.MailNotification,
	teamRepo repository.Team,
) *EmailService {

	isDevelopmentEnv = cfg.IsDevelopmentEnv

	return &EmailService{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Username: cfg.SMTPUser,
		Password: cfg.SMTPPassword,
		From:     cfg.EmailFromEmail,
		FromName: cfg.EmailFromName,
		UseTLS:   cfg.UseMailTLS,

		teamRepo: teamRepo,
	}
}

func (es *EmailService) Notify(bReq *dto.BroadcastSendRequest) (*dto.BroadcastResponse, error) {
	fanEntities, err := es.teamRepo.FindFansByTeamName(context.Background(), bReq.Team)
	if err != nil {
		fmt.Print(err)
	}

	for _, fan := range fanEntities {

		// body := fmt.Sprintf("Time: %s, Tipo: %s, Placar: %s, Mensagem: %s", fan.Team, bReq.Type, bReq.Score, bReq.Message)
		// err := es.rabbitMQ.Publish("broadcast_queue", body)
		// if err != nil {
		// 	log.Println("não foi possível publicar a mensagem")
		// } else {
		// 	log.Println("FOI PRA FILA MANO")
		// }

		es.sendEmail(
			fan.Email,
			fmt.Sprintf("%s da partida do: %s", bReq.Type, bReq.Team),
			bReq.Message,
		)
	}

	return &dto.BroadcastResponse{
		Message: "processando notificações",
	}, nil
}

func (es *EmailService) sendEmail(to string, subject string, body string) error {
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
