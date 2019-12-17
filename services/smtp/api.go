package smtp

import (
	"context"

	"github.com/pkg/errors"
)

func (s *Service) SendEmail(ctx context.Context, email Email) error {
	// setup headers
	headers := make(map[string]string, len(email.Headers)+3)
	headers["From"] = email.From.String()
	headers["Subject"] = email.Subject

	{
		bufTo := s.bufPool.Get()

		encodeAddrs(email.To, bufTo)
		headers["To"] = bufTo.String()

		s.bufPool.Put(bufTo)
	}

	// setup message
	bufBody := s.bufPool.Get()

	for k, v := range headers {
		bufBody.WriteString(k)
		bufBody.Write([]byte(": "))
		bufBody.WriteString(v)
		bufBody.Write([]byte("\r\n"))
	}

	bufBody.Write([]byte("\r\n"))
	bufBody.Write(email.Body)

	client, err := s.getClient()
	if err != nil {
		return errors.Wrap(err, "cannot create new smtp client")
	}

	if err := client.Mail(email.From.Address); err != nil {
		return errors.Wrap(err, "cannot do mail command")
	}

	for i := range email.To {
		if err := client.Rcpt(email.To[i].Address); err != nil {
			return errors.Wrap(err, "cannot do rcpt command for email: "+email.To[i].Address)
		}
	}

	w, err := client.Data()
	if err != nil {
		return errors.Wrap(err, "cannot do data command")
	}

	if _, err := w.Write(bufBody.Bytes()); err != nil {
		return errors.Wrap(err, "cannot write email body")
	}

	s.bufPool.Put(bufBody)

	if err := w.Close(); err != nil {
		return errors.Wrap(err, "cannot close body writer")
	}

	return nil
}
