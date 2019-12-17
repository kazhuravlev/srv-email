package sender

import (
	"net/mail"

	"github.com/kazhuravlev/srv-email/contracts"
	"github.com/kazhuravlev/srv-email/services/smtp"
)

func adaptTaskToEmail(m contracts.Msg) smtp.Email {
	to := make([]mail.Address, len(m.Recipients))
	for i := range m.Recipients {
		to[i] = mail.Address{Address: m.Recipients[i]}
	}

	return smtp.Email{
		From:    mail.Address{Address: m.From},
		To:      to,
		Subject: m.Subject,
		Headers: m.Headers,
		Body:    m.Body,
	}
}
