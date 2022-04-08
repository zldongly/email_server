package service

import (
	"context"
	"go.uber.org/zap"
)

type Template struct {
	Id      string
	Name    string
	Subject string
	Content string
}

type TemplateRepo interface {
	CreateTemplate(ctx context.Context, t *Template) (*Template, error)
}

func NewTempleCase(repo TemplateRepo, log *zap.SugaredLogger) *TemplateUseCase {
	return &TemplateUseCase{
		repo: repo,
		log:  log.With("module", "service/template"),
	}
}

type TemplateUseCase struct {
	repo TemplateRepo
	log  *zap.SugaredLogger
}

func (c *TemplateUseCase) CreateTemplate(ctx context.Context, t *Template) (string, error) {
	temp, err := c.repo.CreateTemplate(ctx, t)
	return temp.Id, err
}
