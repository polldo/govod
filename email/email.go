package email

import (
	"net/smtp"
)

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
	msg := []byte(token)
	dst := []string{to}
	return smtp.SendMail(e.host, e.auth, e.from, dst, msg)
}
