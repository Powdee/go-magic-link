package services

import (
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	dialer *gomail.Dialer
	from   string
}

func NewEmailService(host string, port int, username, password, from string) *EmailService {
	dialer := gomail.NewDialer(host, port, username, password)
	return &EmailService{
		dialer: dialer,
		from:   from,
	}
}

func (es *EmailService) SendMagicLink(email, link string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", es.from)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Your Magic Login Link")
	m.SetBody("text/html", "Click <a href=\""+link+"\">here</a> to log in.")

	if err := es.dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
