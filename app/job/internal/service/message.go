package service

import (
	"context"
	"fmt"
	"github.com/zldongly/email_server/pkg/errors"
	"go.uber.org/zap"
	"strings"
	"time"
)

type MessageRepo interface {
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

type Sender interface {
	Send(receivers []string, subject, content string) error
}

func NewMessageUseCase(repo MessageRepo, s Sender, log *zap.SugaredLogger) *MessageUseCase {
	return &MessageUseCase{
		repo: repo,
		log:  log.With("module", "service"),
		s:   s,
	}
}

type MessageUseCase struct {
	repo MessageRepo
	log  *zap.SugaredLogger

	s Sender
}

func (c *MessageUseCase) SendMail(ctx context.Context, m *Mail) error {
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

	err = c.s.Send(m.Receivers, subject, content)

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
