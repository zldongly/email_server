package api

import (
	"context"
	"github.com/google/wire"
	"github.com/zldongly/email_server/app/admin/internal/service"
	"go.uber.org/zap"
)

var ProviderSet = wire.NewSet(NewAdmin)

type Admin interface {
	CreateTemplate(ctx context.Context, req *CreateTemplateReq) (*CreateTemplateReply, error)
}

type CreateTemplateReq struct {
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

type CreateTemplateReply struct {
	Id string `json:"id"`
}

func NewAdmin(log *zap.SugaredLogger, uc *service.TemplateUseCase) (Admin, error) {
	return &instance{
		log: log.With("module", "api"),
		uc:  uc,
	}, nil
}

type instance struct {
	log *zap.SugaredLogger
	uc  *service.TemplateUseCase
}

func (i *instance) CreateTemplate(ctx context.Context, req *CreateTemplateReq) (*CreateTemplateReply, error) {
	t := &service.Template{
		Name:    req.Name,
		Subject: req.Subject,
		Content: req.Content,
	}

	id, err := i.uc.CreateTemplate(ctx, t)
	if err != nil {
		return nil, err
	}

	return &CreateTemplateReply{
		Id: id,
	}, nil
}
