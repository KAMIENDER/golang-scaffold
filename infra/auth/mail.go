package auth

import (
	"context"

	"github.com/pkg/errors"
	"github.com/volatiletech/authboss/v3"
	"gopkg.in/gomail.v2"
)

var (
	username = ""
	password = ""
)

var _ authboss.Mailer = &EmailSender{}

type EmailSender struct {
}

func (s *EmailSender) Send(ctx context.Context, mail authboss.Email) error {
	if len(mail.To) == 0 {
		return errors.New("[EmailSender] Send len == 0")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", username)
	m.SetHeader("To", mail.To...)
	// m.SetAddressHeader("Cc", "xxx@163.com", "Dan")
	m.SetHeader("Subject", mail.Subject)
	m.SetBody("text/html", mail.TextBody)
	// m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtp.126.com", 25, username, password)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
