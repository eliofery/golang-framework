package email

import (
	"github.com/go-mail/mail/v2"
	"os"
	"strconv"
)

type Email struct {
	From      string
	To        string
	Subject   string
	Plaintext string
	HTML      string
}

type Service struct {
	dialer *mail.Dialer
}

func New() (*Service, error) {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return &Service{}, err
	}

	return &Service{
		dialer: &mail.Dialer{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     port,
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
		},
	}, nil
}

func (es *Service) Send(email Email) error {
	msg := mail.NewMessage()

	msg.SetHeader("From", email.From)
	msg.SetHeader("To", email.To)
	msg.SetHeader("Subject", email.Subject)

	switch {
	case email.Plaintext != "" && email.HTML != "":
		msg.SetBody("text/plain", email.Plaintext)
		msg.AddAlternative("text/html", email.HTML)
	case email.Plaintext != "":
		msg.SetBody("text/plain", email.Plaintext)
	case email.HTML != "":
		msg.AddAlternative("text/html", email.HTML)
	}

	if err := es.dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
