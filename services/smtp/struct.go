package smtp

import "net/mail"

type Email struct {
	From    mail.Address
	To      []mail.Address
	Subject string
	Headers map[string]string
	Body    []byte
}
