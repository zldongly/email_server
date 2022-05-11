package api

import (
	"context"
	"github.com/google/wire"
	"github.com/zldongly/email_server/app/job/internal/service"
	"go.uber.org/zap"
)

var ProviderSet = wire.NewSet(NewJob)

type Job interface {
	SendEmail(ctx context.Context, e *Email) error
}

// Email from kafka msg.value
type Email struct {
	Receivers  []string          `json:"receivers"`
	TemplateId string            `json:"template_id"`
	Param      map[string]string `json:"param"`
}

func NewJob(uc *service.MessageUseCase, log *zap.SugaredLogger) Job {
	return &instance{
		log: log.With("module", "api"),
		uc:  uc,
	}
}

type instance struct {
	log *zap.SugaredLogger
	uc  *service.MessageUseCase
}

func (i *instance) SendEmail(ctx context.Context, e *Email) error {
	m := &service.Mail{
		Receivers:  e.Receivers,
		TemplateId: e.TemplateId,
		Param:      e.Param,
	}
	return i.uc.SendMail(ctx, m)
}
