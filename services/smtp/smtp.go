package smtp

import (
	"crypto/tls"
	"net"
	"net/smtp"

	"github.com/pkg/errors"

	"github.com/kazhuravlev/srv-email/sdk/pool"
)

type Service struct {
	bufPool *pool.BytesBuffer

	host     string
	username string
	password string
	addr     string
}

func New(addr string, username, password string) (*Service, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, errors.Wrap(err, "invalid addr")
	}

	return &Service{
		addr:     addr,
		host:     host,
		username: username,
		password: password,
		bufPool:  pool.NewBytesBuffer(),
	}, nil
}

func (s *Service) getClient() (*smtp.Client, error) {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	// TLS config
	tlsconfig := &tls.Config{
		ServerName: s.host,
	}

	conn, err := tls.Dial("tcp", s.addr, tlsconfig)
	if err != nil {
		return nil, errors.Wrap(err, "cannot dial smtp server")
	}

	smtpClient, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create smtp client")
	}

	if err := smtpClient.Noop(); err != nil {
		return nil, errors.Wrap(err, "cannot do noop")
	}

	// Auth
	if err = smtpClient.Auth(auth); err != nil {
		return nil, errors.Wrap(err, "cannot auth")
	}

	return smtpClient, nil
}
