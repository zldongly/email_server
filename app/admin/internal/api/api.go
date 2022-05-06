package api

import (
	"context"
	"github.com/google/wire"
	"github.com/zldongly/email_server/app/admin/internal/service"
	"go.uber.org/zap"
)

var ProviderSet = wire.NewSet(NewAdmin)

func NewAdmin(log *zap.SugaredLogger, tempUC *service.TemplateUseCase, recordUC *service.RecordUseCase) (Admin, error) {
	return &instance{
		log:      log.With("module", "api"),
		tempUC:   tempUC,
		recordUC: recordUC,
	}, nil
}

type Admin interface {
	CreateTemplate(ctx context.Context, req *CreateTemplateReq) (*CreateTemplateReply, error)
	ListRecord(ctx context.Context, req *ListRecordReq) (*ListRecordReply, error)
}

type CreateTemplateReq struct {
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

type CreateTemplateReply struct {
	Id string `json:"id"`
}

type ListRecordReq struct {
	TemplateId string `json:"template_id" form:"template_id"`
	Name       string `json:"name" form:"name"`
	IsSuccess  *int   `json:"is_success" form:"is_success"`
	PageNum    int    `json:"page_num" form:"page_num"`
	PageSize   int    `json:"page_size" form:"page_size"`
}

type ListRecordReply struct {
	List []*RecordReply `json:"list"`
}

type RecordReply struct {
	Id         string   `json:"id"`
	SendTime   int64    `json:"send_time"`
	Receivers  []string `json:"receivers"`
	TemplateId string   `json:"template_id"`
	Name       string   `json:"name"`
	Content    string   `json:"content"`
	IsSuccess  int8     `json:"is_success"`
	Message    string   `json:"message"`
}

type instance struct {
	log      *zap.SugaredLogger
	tempUC   *service.TemplateUseCase
	recordUC *service.RecordUseCase
}

func (i *instance) CreateTemplate(ctx context.Context, req *CreateTemplateReq) (*CreateTemplateReply, error) {
	t := &service.Template{
		Name:    req.Name,
		Subject: req.Subject,
		Content: req.Content,
	}

	id, err := i.tempUC.CreateTemplate(ctx, t)
	if err != nil {
		return nil, err
	}

	return &CreateTemplateReply{
		Id: id,
	}, nil
}

func (i *instance) ListRecord(ctx context.Context, req *ListRecordReq) (*ListRecordReply, error) {
	list, err := i.recordUC.ListRecord(ctx, req.TemplateId, req.Name, req.IsSuccess, req.PageNum, req.PageSize)
	if err != nil {
		return nil, err
	}

	reply := &ListRecordReply{
		List: make([]*RecordReply, 0, len(list)),
	}
	for _, r := range list {
		reply.List = append(reply.List, &RecordReply{
			Id:         r.Id,
			SendTime:   r.SendTime,
			Receivers:  r.Receivers,
			TemplateId: r.TemplateId,
			Name:       r.Name,
			Content:    r.Content,
			IsSuccess:  r.IsSuccess,
			Message:    r.Message,
		})
	}

	return reply, nil
}
