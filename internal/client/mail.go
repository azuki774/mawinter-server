package client

import (
	"context"
	"fmt"
	"net/smtp"
)

type MailClient struct {
	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
}

func (m *MailClient) Send(ctx context.Context, to string, title string, body string) (err error) {
	recipients := []string{to}

	from := m.SMTPUser
	auth := smtp.PlainAuth("", m.SMTPUser, m.SMTPPass, m.SMTPHost)

	// 送信先は１つのみ対応
	msg := []byte("To: " + to + "\r\n" + "Subject:" + title + "\r\n" + "\r\n" + body)
	if err := smtp.SendMail(fmt.Sprintf("%s:%s", m.SMTPHost, m.SMTPPort), auth, from, recipients, msg); err != nil {
		return err
	}

	return nil
}
