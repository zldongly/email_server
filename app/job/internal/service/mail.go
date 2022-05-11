package service

import (
	"crypto/tls"
	"github.com/zldongly/email_server/app/job/internal/conf"
	"github.com/zldongly/email_server/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"net/http"
)

func NewMailUseCase(cfg *conf.Mail, log *zap.SugaredLogger) Sender {
	log = log.With("module", "service/gomail")
	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	sc, err := d.Dial()
	if err != nil {
		log.Fatal(err)
	}

	return &MailUseCase{
		client: sc,
		cfg:    cfg,
	}
}

type MailUseCase struct {
	client gomail.SendCloser
	cfg    *conf.Mail
}

func (m *MailUseCase) Send(receivers []string, subject, content string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", content)

	if err := m.client.Send(m.cfg.Username, receivers, msg); err != nil {
		err = errors.WithCode(errors.WithStack(err), http.StatusInternalServerError, "mail send")
		return err
	}

	return nil
}
