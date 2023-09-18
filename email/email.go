package email

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"html/template"
	"net/smtp"
)

//go:embed templates
var templates embed.FS

type Emailer struct {
	auth smtp.Auth
	from string
	host string
}

func New(address string, password string, host string, port string) *Emailer {
	a := smtp.PlainAuth("", address, password, host)
	return &Emailer{auth: a, host: host + ":" + port, from: address}
}

func (e *Emailer) SendToken(scope string, token string, to string) error {
	t, err := template.New("email").ParseFS(templates, "templates/reset-password.tmpl")
	if err != nil {
		return fmt.Errorf("parsing email template: %w", err)
	}

	var data struct {
		Link string
	}
	data.Link = "www.google.it"

	var body bytes.Buffer
	err = t.ExecuteTemplate(&body, "html", data)
	if err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	dst := []string{to}
	return smtp.SendMail(e.host, e.auth, e.from, dst, body.Bytes())
}
