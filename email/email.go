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
	auth  smtp.Auth
	from  string
	host  string
	links Links
}

type Links struct {
	ResetPasswordURL string
	ActivationURL    string
}

func New(address string, password string, host string, port string, links Links) *Emailer {
	a := smtp.PlainAuth("", address, password, host)
	return &Emailer{auth: a, host: host + ":" + port, from: address, links: links}
}

func (e *Emailer) SendActivationToken(token string, to string) error {
	t, err := template.New("email").ParseFS(templates, "templates/activation.tmpl")
	if err != nil {
		return fmt.Errorf("parsing email template: %w", err)
	}

	var data struct {
		Link string
	}
	data.Link = e.links.ActivationURL + token

	var body bytes.Buffer
	err = t.ExecuteTemplate(&body, "html", data)
	if err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: Welcome to Govod!\n"
	src := fmt.Sprintf("From: %s\r\n", e.from)
	dst := fmt.Sprintf("To: %s\r\n", to)
	bytes := append([]byte(src+dst+subject+mime), body.Bytes()...)

	return smtp.SendMail(e.host, e.auth, e.from, []string{to}, bytes)
}

func (e *Emailer) SendResetToken(token string, to string) error {
	t, err := template.New("email").ParseFS(templates, "templates/reset-password.tmpl")
	if err != nil {
		return fmt.Errorf("parsing email template: %w", err)
	}

	var data struct {
		Link string
	}
	data.Link = e.links.ResetPasswordURL + token

	var body bytes.Buffer
	err = t.ExecuteTemplate(&body, "html", data)
	if err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: Reset your password\n"
	src := fmt.Sprintf("From: %s\r\n", e.from)
	dst := fmt.Sprintf("To: %s\r\n", to)
	bytes := append([]byte(src+dst+subject+mime), body.Bytes()...)

	return smtp.SendMail(e.host, e.auth, e.from, []string{to}, bytes)
}
