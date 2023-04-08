package auth

import (
	"context"

	"github.com/KAMIENDER/golang-scaffold/infra/config"
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
	host     string
	port     int
	userName string
	passWord string
}

func NewEmailSender(conf *config.Config) *EmailSender {
	return &EmailSender{
		host:     conf.EmailConf.Host,
		port:     conf.EmailConf.Port,
		userName: conf.EmailConf.UserName,
		passWord: conf.EmailConf.PassWord,
	}
}

func (s *EmailSender) Send(ctx context.Context, mail authboss.Email) error {
	if len(mail.To) == 0 {
		return errors.New("[EmailSender] Send len == 0")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", mail.FromName)
	m.SetHeader("To", mail.To...)
	// m.SetAddressHeader("Cc", "xxx@163.com", "Dan")
	m.SetHeader("Subject", mail.Subject)
	m.SetBody("text/html", mail.TextBody)
	// m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer(s.host, s.port, s.userName, s.passWord)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
