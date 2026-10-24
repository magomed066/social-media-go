package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"time"

	"github.com/wneessen/go-mail"
)

const UserWelcomeTemplate = "user_invitation.templ"

//go:embed templates
var FS embed.FS

type WnessenMailer struct {
	from   string
	client *mail.Client
}

func NewWnessenMailer(host, from, username, password string) (*WnessenMailer, error) {
	client, err := mail.NewClient(
		host,
		mail.WithPort(465),
		mail.WithSSL(),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("create SMTP client: %w", err)
	}

	return &WnessenMailer{
		client: client,
		from:   from,
	}, nil
}

func (m *WnessenMailer) Send(to, templateFile string, data any) error {
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return fmt.Errorf("parse email template: %w", err)
	}

	var subject, body bytes.Buffer

	if err := tmpl.ExecuteTemplate(&subject, "subject", data); err != nil {
		return fmt.Errorf("render email subject: %w", err)
	}

	if err := tmpl.ExecuteTemplate(&body, "body", data); err != nil {
		return fmt.Errorf("render email body: %w", err)
	}

	message := mail.NewMsg()

	if err := message.From(m.from); err != nil {
		return fmt.Errorf("set sender: %w", err)
	}

	if err := message.To(to); err != nil {
		return fmt.Errorf("set recipient: %w", err)
	}

	message.Subject(subject.String())
	message.SetBodyString(mail.TypeTextHTML, body.String())

	if err := m.client.DialAndSend(message); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}
