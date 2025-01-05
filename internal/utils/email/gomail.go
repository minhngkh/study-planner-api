package email

import (
	"os"

	"gopkg.in/gomail.v2"
)

var (
	mailHost     = os.Getenv("MAIL_HOST")
	mailAddress  = os.Getenv("MAIL_ADDRESS")
	mailPassword = os.Getenv("MAIL_PASSWORD")
)

func Send(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mailAddress)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(mailHost, 587, mailAddress, mailPassword)

	err := d.DialAndSend(m)
	if err != nil {
		return err
	}

	return nil
}
