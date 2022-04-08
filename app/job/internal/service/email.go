package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/zldongly/email_server/app/job/internal/conf"
	"github.com/zldongly/email_server/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"net/http"
	"strings"
	"time"
)

type EmailRepo interface {
	CreateRecord(ctx context.Context, r *Record) (*Record, error)
	GetTemplate(ctx context.Context, id string) (*Template, error)
}

type Record struct {
	Id         string
	SendTime   int64    // 发送时间
	Receivers  []string // 邮件接受者
	TemplateId string
	Name       string // 模板名称
	Content    string // 邮件内容包括 主题
	IsSuccess  int8   // 1成功 0失败
	Message    string // 失败错误原因
}

type Template struct {
	Id      string
	Name    string
	Subject string
	Content string
}

type Mail struct {
	Receivers  []string
	TemplateId string
	Param      map[string]string
}

func NewEmailUseCase(repo EmailRepo, mc *MailCase, log *zap.SugaredLogger) *EmailUseCase {
	return &EmailUseCase{
		repo: repo,
		log:  log.With("module", "service"),
		mc:   mc,
	}
}

type EmailUseCase struct {
	repo EmailRepo
	log  *zap.SugaredLogger

	mc *MailCase
}

func (c *EmailUseCase) SendMail(ctx context.Context, m *Mail) error {
	t, err := c.repo.GetTemplate(ctx, m.TemplateId)
	if err != nil {
		return err
	}

	c.log.Debug("template:", t)

	// 解析模板
	subject, content, err := t.parse(m)
	if err != nil {
		return err
	}

	err = c.mc.send(m.Receivers, subject, content)

	record := &Record{
		SendTime:   time.Now().Unix(),
		Receivers:  m.Receivers,
		TemplateId: t.Id,
		Name:       t.Name,
		Content:    fmt.Sprintf(`{"subject": "%s", "content"": "%s"}`, subject, content),
		IsSuccess:  1,
		Message:    "",
	}
	if err != nil {
		record.IsSuccess = 0
		record.Message = errors.Cause(err).Error()
	}

	if _, err = c.repo.CreateRecord(ctx, record); err != nil {
		return err
	}

	return nil
}

func (t *Template) parse(m *Mail) (string, string, error) {
	var (
		subject = t.Subject
		content = t.Content
	)

	for k, v := range m.Param {
		subject = strings.ReplaceAll(subject, fmt.Sprintf("{$%s}", k), v)
		content = strings.ReplaceAll(content, fmt.Sprintf("{$%s}", k), v)
	}

	return subject, content, nil
}

func NewMailCase(cfg *conf.Mail, log *zap.SugaredLogger) *MailCase {
	log = log.With("module", "service/gomail")
	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	sc, err := d.Dial()
	if err != nil {
		log.Fatal(err)
	}

	return &MailCase{
		client: sc,
		cfg:    cfg,
	}
}

type MailCase struct {
	client gomail.SendCloser
	cfg    *conf.Mail
}

func (m *MailCase) send(receivers []string, subject, content string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", content)

	if err := m.client.Send(m.cfg.Username, receivers, msg); err != nil {
		err = errors.WithCode(errors.WithStack(err), http.StatusInternalServerError, "mail send")
		return err
	}

	return nil
}
