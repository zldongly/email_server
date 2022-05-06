package service

import (
	"context"
	"go.uber.org/zap"
)

type RecordRepo interface {
	ListRecord(ctx context.Context, templateId, name string, isSuccess *int, pageNum, pageSize int) ([]*Record, error)
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

type RecordUseCase struct {
	repo RecordRepo
	log  *zap.SugaredLogger
}

func NewRecordUseCase(repo RecordRepo, log *zap.SugaredLogger) *RecordUseCase {
	return &RecordUseCase{
		repo: repo,
		log:  log.With("module", "service/record"),
	}
}

func (uc *RecordUseCase) ListRecord(ctx context.Context, templateId, name string, isSuccess *int, pageNum, pageSize int) ([]*Record, error) {
	return uc.repo.ListRecord(ctx, templateId, name, isSuccess, pageNum, pageSize)
}
